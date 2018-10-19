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

//type HashHostsMap struct{
//	hashHostsMap map[string]string
//}

func (m *MsgMapper) AddMsgToMap(msg *nats.Msg) error {

	N := 2 //TODO 共识参数  按validator power

	event := types.Event{}

	if amino.UnmarshalBinary(msg.Data, &event) != nil {

		return errors.New("the event Unmarshal error")
	}

	m.mtx.Lock()
	hashNode, ok := m.MsgMap[event.Sequence]

	//还没有sequence对应记录
	if !ok || hashNode == nil {

		hashNode = make(map[string]string)

		hashNode[string(event.HashBytes)] = event.NodeAddress

		m.MsgMap[event.Sequence] = hashNode

		m.mtx.Unlock()

		return nil
	}

	//sequence已经存在
	if nodes, _ := hashNode[string(event.HashBytes)]; nodes != "" {

		hashNode[string(event.HashBytes)] += "," + event.NodeAddress

		nodes := hashNode[string(event.HashBytes)]

		if strings.Count(nodes, ",") >= N-1 { //达成共识

			log.Infof("consensus from [%s] to [%s] sequence [#%d] hash %s", event.From, event.To, event.Sequence, string(event.HashBytes))
			go m.ferry(event.From, event.To, string(event.HashBytes), nodes, event.Sequence)
		}
	} else {
		hashNode[string(event.HashBytes)] += event.NodeAddress
	}

	m.mtx.Unlock()
	return nil
}
