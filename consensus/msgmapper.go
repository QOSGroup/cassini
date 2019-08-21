package consensus

import (
	"fmt"
	"strings"
	"sync"

	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/pkg/errors"
)

type EngineMap struct {
	mtxMsg sync.RWMutex
	MsgMap map[int64]map[string]string
}

func (m *EngineMap) AddMsgToMap(f *Ferry, event types.Event, N int) (sequence int64, err error) {

	m.mtxMsg.Lock()
	defer m.mtxMsg.Unlock()

	// 关闭共识，收到第一份即共识
	if !f.conf.Consensus {
		h := common.Bytes2HexStr(event.HashBytes)
		n := f.conf.GetQscConfig(event.From).Nodes
		err = f.ConsMap.AddConsToMap(event.Sequence, h, n)
		if err != nil {
			log.Errorf("duplicate AddConsToMap. f.t.s[%s %s #%d] hash [%s]", event.From, event.To, event.Sequence, h[:10])
		}
		delete(m.MsgMap, event.Sequence)
		return event.Sequence + 1, nil
	}
	//----------------

	hashNode, ok := m.MsgMap[event.Sequence]

	////人造拜占庭
	//if strings.Contains(event.NodeAddress, "22") {
	//	log.Infof("add a bai zhan ting sequence #%d", event.Sequence)
	//	event.HashBytes = []byte("baizhanting")
	//}

	log.Debugf("from[%s] seq[%d] ip[%s] hash[%s]", f.from, event.Sequence, event.NodeAddress, event.HashBytes[:10])
	//还没有sequence对应记录
	if !ok || hashNode == nil {

		hashNode = make(map[string]string)

		hashNode[string(event.HashBytes)] = event.NodeAddress

		log.Debugf("msgmapper.AddMsgToMap has no sequence map yet!")
	} else {

		nodes, _ := hashNode[string(event.HashBytes)]
		if !strings.Contains(nodes, event.NodeAddress) {
			if nodes == "" {
				hashNode[string(event.HashBytes)] = event.NodeAddress
			} else {
				hashNode[string(event.HashBytes)] = hashNode[string(event.HashBytes)] + "\000" + event.NodeAddress
			}
		}
	}
	m.MsgMap[event.Sequence] = hashNode

	nodes := hashNode[string(event.HashBytes)]

	if strings.Count(nodes, "\000") >= N-1 {

		hash := common.Bytes2HexStr(event.HashBytes)
		log.Infof("consensus f.t.s[%s %s #%d] hash[%s]", event.From, event.To, event.Sequence, hash[:10])

		//err = f.ferryQCP(event.From, event.To, hash, nodes, event.Sequence)
		nodess := strings.Replace(nodes, "\000", ",", -1)
		err = f.ConsMap.AddConsToMap(event.Sequence, hash, nodess)
		if err != nil {
			log.Errorf("duplicate AddConsToMap. f.t.s[%s %s #%d] hash[%s]", event.From, event.To, event.Sequence, hash[:10])
		}
		delete(m.MsgMap, event.Sequence)
		log.Infof("add Consensus To Map. f.t.s[%s %s #%d] hash[%s]", event.From, event.To, event.Sequence, hash[:10])

		return event.Sequence + 1, nil
	}

	return 0, nil
}

type ConsensusMap struct {
	mtxCon  sync.RWMutex
	ConsMap map[int64]map[string]string
}

func (c *ConsensusMap) AddConsToMap(sequence int64, hash, nodes string) error {

	c.mtxCon.Lock()
	defer c.mtxCon.Unlock()

	hashNode, ok := c.ConsMap[sequence]

	if !ok || hashNode == nil {

		hashNode = make(map[string]string)

		hashNode[hash] = nodes

		c.ConsMap[sequence] = hashNode

	} else {
		return errors.New("duplicate AddConsToMap")
	}
	return nil
}

type Consensus struct {
	Sequence int64
	Hash     string
	Nodes    string
}

func (c *ConsensusMap) GetConsFromMap(sequence int64) (*Consensus, error) {
	c.mtxCon.Lock()
	defer c.mtxCon.Unlock()

	cons := Consensus{}
	cons.Sequence = sequence

	hashNode, ok := c.ConsMap[sequence]

	if !ok || hashNode == nil {
		return nil, fmt.Errorf("not found consensus, sequence: %d", sequence)
	}
	for k, v := range hashNode {
		cons.Hash = k
		cons.Nodes = v
	}
	return &cons, nil
}
