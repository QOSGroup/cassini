package consensus

import (
	"strings"
	"sync"

	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
)

type MsgMapper struct {
	mtx    sync.RWMutex
	MsgMap map[int64]map[string]string
}

func (m *MsgMapper) AddMsgToMap(f *Ferry, event types.Event, N int) (sequence int64, err error) {

	m.mtx.Lock()
	defer m.mtx.Unlock()

	// 仅为测试，临时添加
	if strings.EqualFold("no", f.conf.Consensus) {
		h := common.Bytes2HexStr(event.HashBytes)
		n := f.conf.GetQscConfig(event.From).NodeAddress
		err = f.ferryQCP(event.From, event.To, h, n, event.Sequence)
		if err != nil {
			return 0, err
		}
		return event.Sequence + 1, nil
	}
	//----------------

	hashNode, ok := m.MsgMap[event.Sequence]

	//log.Infof("%v", m.MsgMap)l

	//还没有sequence对应记录
	if !ok || hashNode == nil {

		hashNode = make(map[string]string)

		hashNode[string(event.HashBytes)] = event.NodeAddress

		m.MsgMap[event.Sequence] = hashNode
		log.Debugf("msgmapper.AddMsgToMap has no sequence map yet!")
	} else {

		nodes, _ := hashNode[string(event.HashBytes)]
		if !strings.Contains(nodes, event.NodeAddress) {
			hashNode[string(event.HashBytes)] += "," + event.NodeAddress
		}
	}

	nodes := hashNode[string(event.HashBytes)]

	if strings.Count(nodes, ",") >= N-1 {

		hash := common.Bytes2HexStr(event.HashBytes)
		log.Infof("consensus from [%s] to [%s] sequence [#%d] hash [%s]", event.From, event.To, event.Sequence, hash[:10])

		err = f.ferryQCP(event.From, event.To, hash, nodes, event.Sequence)
		if err != nil {
			return 0, err
		}
		delete(m.MsgMap, event.Sequence)
		log.Debugf("msgmapper.AddMsgToMap ferryQCP")
		return event.Sequence + 1, nil

	}

	return 0, nil
}
