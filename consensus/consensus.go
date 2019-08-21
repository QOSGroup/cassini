package consensus

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/queue"
	"github.com/QOSGroup/cassini/types"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"

	"github.com/QOSGroup/cassini/restclient"
	"github.com/tendermint/tendermint/libs/common"
)

type consResult int

const (
	fail consResult = iota
	success
	empty //交易还没产生 共识失败
)

// ConsEngine Consensus engine
type ConsEngine struct {
	M        *EngineMap
	F        *Ferry
	sequence int64
	from     string
	to       string
	mtx      sync.RWMutex
}

var wg sync.WaitGroup

// StartQcpConsume Start to consume tx msg
func StartQcpConsume(conf *config.Config) (err error) {

	// qsconfigs := config.DefaultQscConfig()
	qsconfigs := conf.Qscs

	if len(qsconfigs) < 2 {
		return errors.New("config error , at least two chain targets ")
	}

	var subjects string

	es := make(chan error, 1024) //TODO 1024参数按需修改
	defer close(es)

	engines := make([]*ConsEngine, 0)
	var ce *ConsEngine
	for i, qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs); j++ {
			wg.Add(2)

			ce = createConsEngine(qsconfigs[j].Name, qsconfig.Name, conf, es)
			engines = append(engines, ce)
			ce = createConsEngine(qsconfig.Name, qsconfigs[j].Name, conf, es)
			engines = append(engines, ce)

			subjects += fmt.Sprintf("[%s] [%s]", qsconfigs[j].Name+"2"+qsconfig.Name, qsconfig.Name+"2"+qsconfigs[j].Name)

		}
	}

	wg.Wait()

	if len(es) > 0 {
		for e := range es {
			log.Error(e)
		}
		return errors.New("couldn't start qcp consumer")
	}

	log.Infof("listening on subjects %s", subjects)

	for _, ce := range engines {
		ce.StartEngine()
		ce.F.StartFerry()
	}

	ticker := func(engines []*ConsEngine) {
		log.Debugf("run roomkeeper...%d", len(engines))
		// 定时触发共识引擎
		tick := time.NewTicker(time.Duration(conf.EventWaitMillitime*10) * time.Millisecond)
		for range tick.C {
			log.Debug("run roomkeeper...")
			for _, ce := range engines {
				ce.RoomKeeper()
			}
		}
	}
	go ticker(engines)

	return
}

func createConsEngine(from, to string, conf *config.Config, e chan<- error) (ce *ConsEngine) {
	ce = NewConsEngine(from, to)

	seq, err := ce.F.GetSequenceFromChain(from, to, "in") // seq= toChain's in/fromchain/maxseq
	if err != nil {
		log.Errorf("Create consensus engine error: %v", err)
	} else {
		log.Debugf("Create consensus engine query chain %s in-sequence: %d", to, seq)
		ce.SetSequence(from, to, seq)
	}

	go qcpConsume(ce, from, to, conf, e)
	return ce
}

//QcpConsumer consume the message from nats server
//
// from ,to is chain name for example "QOS"
func qcpConsume(ce *ConsEngine, from, to string, conf *config.Config, e chan<- error) {
	log.Debugf("Consume qcp f.t[%s %s]", from, to)

	var i int64

	defer wg.Done()

	listener := func(data []byte, consumer queue.Consumer) {
		i++

		tx := types.Event{}
		amino.UnmarshalBinaryLengthPrefixed(data, &tx)

		log.Infof("[#%d] Consume subject [%s] sequence [#%d] nodeAddress '%s'",
			i, consumer.Subject(), tx.Sequence, tx.NodeAddress)

		// 监听到交易事件后立即查询需要等待一段时间才能查询到交易数据；
		//TODO 优化
		// 需要监听下一个块的New Block 事件以确认交易数据入块，abci query 接口才能够查询出交易；
		// 同时提供定时触发机制，以保证共识模块在交易事件丢失或网络错误等问题出现时仍然能够正常运行。
		//if conf.EventWaitMillitime > 0 {
		//	time.Sleep(time.Duration(conf.EventWaitMillitime) * time.Millisecond)
		//}

		ce.Add2Engine(data)
	}

	// consummer := msgqueue.NATSConsumer{
	// 	ServerUrls: conf.Queue,
	// 	Subject:    from + "2" + to,
	// 	CallBack:   cb}

	// nc, err := consummer.Connect()
	// if err != nil {
	// 	e <- err
	// 	return
	// }

	// if err = consummer.Consume(nc); err != nil {
	// 	e <- err
	// 	return
	// }

	subject := from + "2" + to

	consumer, err := queue.NewConsumer(subject)
	if err != nil {
		e <- err
	}
	if consumer == nil {
		e <- fmt.Errorf("New consumer error: get nil")
		return
	}
	consumer.Subscribe(listener)
	return
}

// NewConsEngine New a consensus engine
func NewConsEngine(from, to string) *ConsEngine {
	ce := new(ConsEngine)
	ce.M = &EngineMap{MsgMap: make(map[int64]map[string]string)}
	ce.F = NewFerry(config.GetConfig(), from, to, 0)

	ce.from = from
	ce.to = to
	return ce
}

// Add2Engine Add a message to consensus engine
func (c *ConsEngine) Add2Engine(data []byte) error {
	event := types.Event{}

	if amino.UnmarshalBinaryLengthPrefixed(data, &event) != nil {

		return errors.New("the event Unmarshal error")
	}

	if event.Sequence < c.sequence {
		return errors.New("msg sequence is small then the sequence in consensus engine")
	}

	_, err := c.M.AddMsgToMap(c.F, event, c.consensus32())
	if err != nil {
		return err
	}
	//c.F.SetSequence(seq)
	return nil
}

func (c *ConsEngine) consensus32() (N int) {
	nodes := c.F.conf.GetQscConfig(c.from).Nodes

	n := len(strings.Split(nodes, ","))

	N = (n*2 + 2) / 3

	log.Debugf("f.t[%s %s] [consensus N #%d]", c.from, c.to, N)
	return int(N)
}

// StartEngine 启动共识引擎尝试处理下一个交易
func (c *ConsEngine) StartEngine() error {
	go func() {
		for {
			_, err := c.F.ConsMap.GetConsFromMap(c.sequence)
			if err == nil { //已有共识
				c.SetSequence(c.from, c.to, c.sequence)
			}

			cresult := c.conSequence()
			if cresult == success {
				c.SetSequence(c.from, c.to, c.sequence)
			}
			if cresult == empty {
				time.Sleep(time.Duration(c.F.conf.EventWaitMillitime) * time.Millisecond)
				continue
			}
			if cresult == fail { //TODO 不能达成共识 继续下一sequence？
				log.Errorf("MsgMap%v", c.M.MsgMap[c.sequence])
				time.Sleep(time.Duration(c.F.conf.EventWaitMillitime) * time.Millisecond * 10)
				//s := fmt.Sprintf("consensusEngine f.t.s[%s %s #%d] failed.", c.from, c.to, c.sequence)
				//panic(s)
			}
		}
	}()

	return nil
}

// conSequence 交易还没产生和共识出错区别开
func (c *ConsEngine) conSequence() consResult {

	log.Debugf("Start consensus engine f.t.s[%s %s #%d]", c.from, c.to, c.sequence)

	nodes := c.F.conf.GetQscConfig(c.from).Nodes

	N := c.consensus32()

	var bempty bool
	for _, node := range strings.Split(nodes, ",") {

		qcp, err := c.F.queryTxQcpFromNode(c.from, c.to, node, c.sequence) // be (c.to, node, c.sequence)

		if qcp == nil {
			if err != nil && strings.Contains(err.Error(), restclient.ERR_emptyqcp) {
				bempty = true //交易还没产生
			}
			continue
		}
		hash := crypto.Sha256(qcp.BuildSignatureBytes())
		ced := types.CassiniEventDataTx{From: c.from, To: c.to, Height: qcp.BlockHeight, Sequence: c.sequence}

		ced.HashBytes = hash

		event := types.Event{NodeAddress: node, CassiniEventDataTx: ced}

		seq, err := c.M.AddMsgToMap(c.F, event, N)
		if err != nil {
			continue
		}
		if seq > 0 {
			return success
		}
	}

	if bempty {
		return empty
	}
	return fail
}

// SetSequence 设置交易序列号
func (c *ConsEngine) SetSequence(from, to string, s int64) {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	seq, _ := c.GetSequenceFromChain(from, to, "in")

	c.sequence = common.MaxInt64(s, seq) + 1
	log.Infof("f.t[%s %s] ConsEngine sequence set to [#%d]", from, to, c.sequence)
}

// GetSequenceFromChain 在to chain上查询 来自/要去 from chain 的 sequence
func (c *ConsEngine) GetSequenceFromChain(from, to, inout string) (int64, error) {
	//qsc := c.F.conf.GetQscConfig(to)
	//
	//nodeto := strings.Split(qsc.NodeAddress, ",")
	//
	//add := GetAddressFromUrl(nodeto[0]) //TODO 多node 取sequence
	//r := c.F.rmap[add]
	//
	//return r.GetSequence(from, inout)
	return c.F.GetSequenceFromChain(from, to, inout)
}

// GetAddress returns address
func GetAddress(url string) string {
	n := strings.Index(url, "://")
	if n < 0 {
		return url
	}
	return url[n+3:]
}

// Setfrom set from value
func (c *ConsEngine) Setfrom(from string) {
	c.from = from
}

// Setto set to value
func (c *ConsEngine) Setto(to string) {
	c.to = to
}

// Getfrom return from value
func (c *ConsEngine) Getfrom() string {
	return c.from
}

// Getto returns to value
func (c *ConsEngine) Getto() string {
	return c.to
}

// RoomKeeper 清理 EngineMap,ConsensusMap过时k-v;校正sequence
func (c *ConsEngine) RoomKeeper() {
	c.cleanMap()
	c.ajustSequence()
}

func (c *ConsEngine) cleanMap() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.F.mtx.Lock()
	defer c.F.mtx.Unlock()

	for k := range c.M.MsgMap {
		if k < c.F.sequence {
			delete(c.M.MsgMap, k)
			log.Debugf("delete c.sequence[#%d]", k)
		}
	}
	for k := range c.F.ConsMap.ConsMap {
		if k < c.F.sequence {
			delete(c.F.ConsMap.ConsMap, k)
			log.Debugf("delete f.sequence[#%d]", k)
		}
	}

}

func (c *ConsEngine) ajustSequence() {

	seqDes, _ := c.F.GetSequenceFromChain(c.from, c.to, "in")
	seqSou, _ := c.F.GetSequenceFromChain(c.to, c.from, "out")

	//if seqDes >= seqSou || c.sequence > seqSou {
	//	time.Sleep(time.Duration(c.F.conf.EventWaitMillitime) * time.Millisecond)
	//	continue
	//}

	if seqSou < seqDes {
		log.Warn("the destnation sequence bigger than source sequence")
		//panic("the destnation sequence bigger than source sequence")
	}
	if seqDes >= c.sequence {
		c.SetSequence(c.from, c.to, seqDes)
	}

	if seqDes >= c.F.sequence {
		c.F.SetSequence(c.F.from, c.F.to, seqDes)
	}

}
