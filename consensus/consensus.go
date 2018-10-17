package consensus

import (
	"errors"
	"strings"
	"sync"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbasebak/txs"
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
			go m.ferry(event.From, event.To, string(event.HashBytes), nodes, event.Sequence)
		}
	} else {
		hashNode[string(event.HashBytes)] += event.NodeAddress
	}

	m.mtx.Unlock()
	return nil
}

//ferry get qcp transaction from source chain and post it to destnation chain
//from is chain name of the source chain
//to is the chain name of destnation chain
//nodes is consensus nodes of the source chain
func (m *MsgMapper) ferry(from, to, hash, nodes string, sequence int64) (err error) {

	qcp := new(txs.TxQcp)
	success := false

EndGet:

	for _, node := range strings.Split(nodes, ",") {
		r := restclient.NewRestClient(node) //"tcp://192.168.168.195:26657"
		qcp, err = r.GetTxQcp(to, sequence)
		if err != nil {
			continue
		}

		//TODO qcp hash 与 hash值比对
		if e := qcp.ValidateBasicData(true, "QOS"); &e == nil {
			success = true
			break EndGet
		}

		//TODO qxtcp联盟链公钥验签 baseapp.validateTxQcpSignature 签名过程basecli\main.go genQcpSendTx

	}

	if !success {
		return errors.New("get qcp transaction failed")
	}

	qscConfig := config.DefaultQscConfig()
	toNodes := qscConfig[0].NodeAddress //TODO 取目标链nodes 地址

EndPost:
	for _, node := range strings.Split(toNodes, ",") {

		r := restclient.NewRestClient(node)
		err := r.PostTxQcp(to, qcp) //TODO 连接每个目标链node
		if err != nil {
			continue
		}

		log.Info("ferried from [%s] to [%s] sequence [#%d] hash %s", qcp.From, to, sequence, hash)

		success = true
		break EndPost
	}

	if success {
		return nil
	}

	return errors.New("ferry qcp transaction failed")

}
