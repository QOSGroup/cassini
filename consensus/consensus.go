package consensus

import (
	"errors"
	"strings"
	"sync"

	"crypto/sha256"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
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
//
//from is chain name of the source chain
//to is the chain name of destnation chain
//nodes is consensus nodes of the source chain
func (m *MsgMapper) ferry(from, to, hash, nodes string, sequence int64) (err error) {

	qcp, err := m.getTxQcp(from, to, hash, nodes, sequence)

	if err != nil {
		return errors.New("get qcp transaction failed")
	}

	//TODO 公链签名
	if from == "QOS" {

	}

	//TODO 取目标链nodes 地址

	err = m.postTxQcp(to, qcp)

	if err != nil {
		return errors.New("post qcp transaction failed")
	}

	log.Infof("success ferry txqcp from [%s] to [%s] sequence [#%d]", from, to, sequence)
	return nil

}

//getTxQcp get QCP transactions from sorce chain
func (m *MsgMapper) getTxQcp(from, to, hash, nodes string, sequence int64) (qcp *txs.TxQcp, err error) {

	success := false

EndGet:

	for _, node := range strings.Split(nodes, ",") {

		qcp, err := m.getTxQcpFromNode(to, hash, node, sequence)

		if err != nil || qcp == nil {
			continue
		}

		success = true
		break EndGet

	}

	if !success {
		return nil, errors.New("get qcp transaction from chain " + from + " failed")
	}

	return
}

func (m *MsgMapper) getTxQcpParalle(from, to, hash, nodes string, sequence int64) (qcps []txs.TxQcp, err error) {

	nodeList := strings.Split(nodes, ",")
	var tasks = make([]common.Task, len(nodeList))

	for i := 0; i < len(tasks); i++ {
		tasks[i] = func(i int) (res interface{}, err error, abort bool) {
			qcp, err := m.getTxQcpFromNode(to, hash, nodeList[i], sequence)
			return qcp, err, false //TODO
		}
	}

	var tResults, ok = common.Parallel(tasks...)
	if !ok {
		log.Error("parallel failed")
	}

	var failTasks int
	for i := 0; i < len(tasks); i++ {
		tResult, ok := tResults.LatestResult(i)
		if !ok {
			failTasks++
		} else if tResult.Error != nil {
			failTasks++
		} else {
			qcps = append(qcps, *(tResult.Value).(*txs.TxQcp))
		}

	}

	if len(qcps)*2 > failTasks { //TODO 加入共识逻辑
		return qcps, nil
	}

	return nil, errors.New("parallel get qcp transaction from chain " + from + " failed")
}

//getTxQcpFromNode get QCP transactions from single chain node
func (m *MsgMapper) getTxQcpFromNode(to, hash, node string, sequence int64) (qcp *txs.TxQcp, err error) {

	r := restclient.NewRestClient(node) //"tcp://127.0.0.1:26657"
	qcp, err = r.GetTxQcp(to, sequence)
	if err != nil || qcp == nil {
		return nil, errors.New("get TxQcp from " + node + "failed.")
	}

	//TODO 取本地联盟链公钥验签
	pubkey := qcp.Sig.Pubkey
	if !pubkey.VerifyBytes(qcp.GetSigData(), qcp.Sig.Signature) {
		return nil, errors.New("get TxQcp from " + node + " data verify failed.")
	}

	//TODO qcp hash 与 hash值比对
	h := sha256.New()
	h.Write(qcp.GetSigData())
	if string(h.Sum(nil)) == hash {
		return qcp, nil
	}

	return nil, errors.New("get TxQcp from " + node + "failed")
}

func (m *MsgMapper) postTxQcp(to string, qcp *txs.TxQcp) (err error) {

	success := false
	qscConfig := config.DefaultQscConfig()
	toNodes := qscConfig[0].NodeAddress //TODO 取目标链nodes 地址

EndPost:
	for _, node := range strings.Split(toNodes, ",") {

		r := restclient.NewRestClient(node)
		err := r.PostTxQcp(to, qcp) //TODO 连接每个目标链node
		if err != nil {
			continue
		}

		success = true
		break EndPost
	}

	if !success {
		return errors.New("post qcp transaction failed")
	}

	return

}
