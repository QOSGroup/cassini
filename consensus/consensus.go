package consensus

import (
	"errors"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"

	"github.com/QOSGroup/cassini/types"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"

	"fmt"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/tendermint/tendermint/libs/common"
	"strings"
	"sync"
	"time"
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
func (c *ConsEngine) Add2Engine(msg *nats.Msg) error {
	event := types.Event{}

	if amino.UnmarshalBinary(msg.Data, &event) != nil {

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
	nodes := c.F.conf.GetQscConfig(c.from).NodeAddress

	n := len(strings.Split(nodes, ","))

	N = (n*2 + 2) / 3

	log.Debugf("f.t[%s %s] [consensus N #%d]", c.from, c.to, N)
	return int(N)
}

// StartEngine 启动共识引擎尝试处理下一个交易
func (c *ConsEngine) StartEngine() error {

	for {
		//seqDes, _ := c.F.GetSequenceFromChain(c.from, c.to, "in")
		//seqSou, _ := c.F.GetSequenceFromChain(c.to, c.from, "out")
		//
		//if seqDes >= seqSou || c.sequence > seqSou {
		//	time.Sleep(time.Duration(c.F.conf.EventWaitMillitime) * time.Millisecond)
		//	continue
		//}
		//
		//if seqDes >= c.sequence {
		//	c.SetSequence(c.from, c.to, seqDes)
		//}

		_, err := c.F.ConsMap.GetConsFromMap(c.sequence)
		if err == nil { //已有共识
			c.SetSequence(c.from, c.to, c.sequence)
		}

		cresult := c.ConSequence()
		if cresult == success {
			c.SetSequence(c.from, c.to, c.sequence)
		}
		if cresult == empty {
			time.Sleep(time.Duration(c.F.conf.EventWaitMillitime) * time.Millisecond)
			continue
		}
		if cresult == fail { //TODO 不能达成共识 继续下一sequence？
			log.Errorf("MsgMap%v", c.M.MsgMap[c.sequence])
			s := fmt.Sprintf("consensusEngine f.t.s[%s %s #%d] failed.", c.from, c.to, c.sequence)
			panic(s)
		}
	}

	return nil
}

func (c *ConsEngine) ConSequence() consResult { //交易还没产生和共识出错区别开

	log.Debugf("Start consensus engine f.t.s[%s %s #%d]", c.from, c.to, c.sequence)

	nodes := c.F.conf.GetQscConfig(c.from).NodeAddress

	N := c.consensus32()

	var bempty bool
	for _, node := range strings.Split(nodes, ",") {

		qcp, err := c.F.queryTxQcpFromNode(c.to, node, c.sequence) // be (c.to, node, c.sequence)

		if err != nil || qcp == nil {
			if strings.Contains(err.Error(), restclient.ERR_emptyqcp) {
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

//在to chain上查询 来自/要去 from chain 的 sequence
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

func GetAddressFromUrl(url string) string {
	n := strings.Index(url, "://")
	if n < 0 {
		return url
	}
	return url[n+3:]
}

func (c *ConsEngine) Setfrom(from string) {
	c.from = from
}

func (c *ConsEngine) Setto(to string) {
	c.to = to
}

func (c *ConsEngine) Getfrom() string {
	return c.from
}

func (c *ConsEngine) Getto() string {
	return c.to
}

//roomKeeper 清理 EngineMap,ConsensusMap过时k-v;校正sequence
func (c *ConsEngine) RoomKeeper() {
	c.cleanMap()
	c.ajustSequence()
}

func (c *ConsEngine) cleanMap() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.F.mtx.Lock()
	defer c.F.mtx.Unlock()

	for k, _ := range c.M.MsgMap {
		if k < c.F.sequence {
			delete(c.M.MsgMap, k)
			log.Debugf("delete c.sequence[#%d]", k)
		}
	}
	for k, _ := range c.F.ConsMap.ConsMap {
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
		panic("the destnation sequence bigger than source sequence")
	}
	if seqDes >= c.sequence {
		c.SetSequence(c.from, c.to, seqDes)
	}

	if seqDes >= c.F.sequence {
		c.F.SetSequence(c.F.from, c.F.to, seqDes)
	}

}
