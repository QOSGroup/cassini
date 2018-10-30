package consensus

import (
	"errors"
	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"strings"

	"github.com/tendermint/tendermint/libs/common"
)

type ConsEngine struct {
	M        *MsgMapper
	f        *Ferry
	sequence int64
}

func NewConsEngine() *ConsEngine {
	ce := new(ConsEngine)
	ce.M = &MsgMapper{MsgMap: make(map[int64]map[string]string)}
	ce.f = &Ferry{}
	return ce
}

func (c *ConsEngine) Add2Engine(msg *nats.Msg) error {
	event := types.Event{}

	if amino.UnmarshalBinary(msg.Data, &event) != nil {

		return errors.New("the event Unmarshal error")
	}

	if event.Sequence < c.sequence {
		return errors.New("msg sequence is small then the sequence in consensus engine")
	}

	seq, err := c.M.AddMsgToMap(event, c.f)
	if err != nil {
		return err
	}
	c.setSequence(seq)
	return nil
}

func (c *ConsEngine) setSequence(s int64) {
	c.sequence = s
}

type Ferry struct {
}

//ferryQCP get qcp transaction from source chain and post it to destnation chain
//
//from is chain name of the source chain
//to is the chain name of destnation chain
//nodes is consensus nodes of the source chain
func (f *Ferry) ferryQCP(from, to, hash, nodes string, sequence int64) (err error) {

	log.Debugf("Ferry qcp from [%s] to [%s], sequence=%d", from, to, sequence)

	qcp, err := f.getTxQcp(from, to, hash, nodes, sequence)

	if err != nil {
		log.Errorf("%v", err)
		return errors.New("get qcp transaction failed")
	}

	//TODO 公链签名
	if from == "QOS" {

	}

	//TODO 取目标链nodes 地址

	err = f.postTxQcp(to, qcp)

	if err != nil {
		return errors.New("post qcp transaction failed")
	}

	log.Infof("success ferry qcp transaction from [%s] to [%s] sequence [#%d] \n", from, to, sequence)
	return nil

}

//getTxQcp get QCP transactions from sorce chain
func (f *Ferry) getTxQcp(from, to, hash, nodes string, sequence int64) (qcp *txs.TxQcp, err error) {

	success := false

EndGet:

	for _, node := range strings.Split(nodes, ",") {

		qcp, err = f.getTxQcpFromNode(to, hash, node, sequence)

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

func (f *Ferry) getTxQcpParalle(from, to, hash, nodes string, sequence int64) (qcps []txs.TxQcp, err error) {

	nodeList := strings.Split(nodes, ",")
	var tasks = make([]common.Task, len(nodeList))

	for i := 0; i < len(tasks); i++ {
		tasks[i] = func(i int) (res interface{}, err error, abort bool) {
			qcp, err := f.getTxQcpFromNode(to, hash, nodeList[i], sequence)
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
func (f *Ferry) getTxQcpFromNode(to, hash, node string, sequence int64) (qcp *txs.TxQcp, err error) {

	r := restclient.NewRestClient(node) //"tcp://127.0.0.1:26657"
	qcp, err = r.GetTxQcp(to, sequence)
	if err != nil || qcp == nil {
		return nil, errors.New("get TxQcp from " + node + "failed.")
	}

	//TODO 取本地联盟链公钥验签
	//pubkey := qcp.Sig.Pubkey  //mock pubkey 为 nil pnic
	//if !pubkey.VerifyBytes(qcp.GetSigData(), qcp.Sig.Signature) {
	//	return nil, errors.New("get TxQcp from " + node + " data verify failed.")
	//}

	//TODO qcp hash 与 hash值比对
	//if string(tmhash.Sum(qcp.GetSigData())) != hash { //算法保持 tmhash.hash 一致 sha256 前 20byte

	hash2 := cmn.Bytes2HexStr(crypto.Sha256(qcp.GetSigData()))
	if hash2 != hash {
		return nil, errors.New("get TxQcp from " + node + "failed")
	}

	return qcp, nil

}

func (f *Ferry) postTxQcp(to string, qcp *txs.TxQcp) (err error) {

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
