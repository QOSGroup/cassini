package consensus

import (
	"errors"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
	"strings"
	"sync"
)

type MsgMapper struct {
	mtx    sync.RWMutex
	MsgMap map[int64]map[string]string
}

func (m *MsgMapper) AddMsgToMap(msg *nats.Msg, f *Ferry) error {

	N := 2 //TODO 共识参数  按validator voting power

	event := types.Event{}

	if amino.UnmarshalBinary(msg.Data, &event) != nil {

		return errors.New("the event Unmarshal error")
	}

	m.mtx.Lock()

	defer m.mtx.Unlock()

	hashNode, ok := m.MsgMap[event.Sequence]

	//log.Infof("%v", m.MsgMap)

	//还没有sequence对应记录
	if !ok || hashNode == nil {

		hashNode = make(map[string]string)

		hashNode[string(event.HashBytes)] = event.NodeAddress

		m.MsgMap[event.Sequence] = hashNode

		return nil
	}

	//sequence已经存在
	if nodes, _ := hashNode[string(event.HashBytes)]; nodes != "" {

		if strings.Contains(nodes, event.NodeAddress) { //有节点重复广播event
			return nil
		}

		hashNode[string(event.HashBytes)] += "," + event.NodeAddress

		nodes := hashNode[string(event.HashBytes)]

		if strings.Count(nodes, ",") >= N-1 { //TODO 达成共识

			log.Infof("consensus from [%s] to [%s] sequence [#%d] hash %s", event.From, event.To, event.Sequence, string(event.HashBytes))

			//f := Ferry{}
			go f.ferryQCP(event.From, event.To, string(event.HashBytes), nodes, event.Sequence)
		}
	} else {

		hashNode[string(event.HashBytes)] += event.NodeAddress
	}

	return nil
}
