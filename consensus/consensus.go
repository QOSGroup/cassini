package consensus

import (
	"errors"
	"strings"
	"sync"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
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
			//go m.ferry(event.From, event.To, string(event.HashBytes), nodes, event.Sequence)
		}
	} else {
		hashNode[string(event.HashBytes)] += event.NodeAddress
	}

	m.mtx.Unlock()
	return nil
}

func (m *MsgMapper) ferry(from, to, hash, nodes string, sequence int64) error {

	for _, node := range strings.Split(nodes, ",") {
		r := restclient.NewRestClient("tcp://" + node) //"tcp://192.168.168.195:26657"
		qcp, err := r.GetTxQcp(to, sequence)
		if err != nil {
			continue
		}

		//TODO qcp hash 与 hash值比对

		//TODO qxtcp联盟链公钥验签

		err = r.PostTxQcp(to, qcp) //TODO 连接每个目标链node
		if err != nil {
			continue
		}

		log.Info("ferried from [%s] to [%s] sequence [#%d] hash %s", qcp.From, to, sequence, hash)

		return nil
	}
	return errors.New("ferry qcptx failed")

}
