package consensus

import (
	"errors"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"

	"github.com/QOSGroup/cassini/types"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/libs/common"
	"strings"
	"sync"
	"time"
)

// ConsEngine Consensus engine
type ConsEngine struct {
	M        *EngineMap
	F        *Ferry
	sequence int64
	from     string
	to       string
	mtx      sync.RWMutex
	//conf     *config.Config
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

	log.Infof("from [%s] to [%s] [consensus N #%d]", c.from, c.to, N)
	return int(N)
}

// StartEngine 启动共识引擎尝试处理下一个交易
func (c *ConsEngine) StartEngine() error {

	for {
		seqDes, _ := c.F.GetSequenceFromChain(c.from, c.to, "in")
		seqSou, _ := c.F.GetSequenceFromChain(c.to, c.from, "out")

		if seqDes >= seqSou || c.sequence > seqSou {
			time.Sleep(time.Duration(c.F.conf.EventWaitMillitime) * time.Millisecond)
			continue
		}

		if seqDes >= c.sequence {
			c.SetSequence(c.from, c.to, seqDes)
		}

		_, err := c.F.ConsMap.GetConsFromMap(c.sequence)
		if err == nil { //已有共识
			c.SetSequence(c.from, c.to, c.sequence)
		}
		if c.ConSequence() {
			c.SetSequence(c.from, c.to, c.sequence)
		} else {
			log.Errorf("consensusEngine from [%s] to [%s] sequence [#%d] failed. seqDes [%d] seqSou[%d] ", c.from, c.to, c.sequence, seqDes, seqSou)
		}
	}

	return nil
}

func (c *ConsEngine) ConSequence() bool {

	log.Debugf("Start consensus engine from: [%s] to: [%s] sequence: [%d]", c.from, c.to, c.sequence)

	nodes := c.F.conf.GetQscConfig(c.from).NodeAddress

	N := c.consensus32()

	for _, node := range strings.Split(nodes, ",") {

		qcp, err := c.F.queryTxQcpFromNode(c.to, node, c.sequence) // be (c.to, node, c.sequence)

		if err != nil || qcp == nil {
			continue
		}
		hash := crypto.Sha256(qcp.GetSigData())
		ced := types.CassiniEventDataTx{From: c.from, To: c.to, Height: qcp.BlockHeight, Sequence: c.sequence}

		ced.HashBytes = hash

		event := types.Event{NodeAddress: node, CassiniEventDataTx: ced}

		seq, err := c.M.AddMsgToMap(c.F, event, N)
		if err != nil {
			continue
		}
		if seq > 0 {
			return true
		}
	}

	return false
}

// SetSequence 设置交易序列号
func (c *ConsEngine) SetSequence(from, to string, s int64) {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	seq, _ := c.GetSequenceFromChain(from, to, "in")

	c.sequence = common.MaxInt64(s, seq) + 1
	log.Infof("from [%s] to [%s] ConsEngine sequence set to [#%d]", from, to, c.sequence)
}

//在to chain上查询 来自/要去 from chain 的 sequence
func (c *ConsEngine) GetSequenceFromChain(from, to, inout string) (int64, error) {
	qsc := c.F.conf.GetQscConfig(to)

	nodeto := strings.Split(qsc.NodeAddress, ",")

	add := GetAddressFromUrl(nodeto[0]) //TODO 多node 取sequence
	r := c.F.rmap[add]

	return r.GetSequence(from, inout)
}

func GetAddressFromUrl(url string) string {
	n := strings.Index(url, "://")
	if n < 0 {
		return url
	}
	return url[n+3:]
}
