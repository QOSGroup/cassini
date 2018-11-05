package consensus

import (
	"errors"
	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/example/basecoin/app"
	"github.com/QOSGroup/qbase/txs"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	"strings"
	"sync"
)

// ConsEngine Consensus engine
type ConsEngine struct {
	M *MsgMapper
	F *Ferry
	//sequence int64
	from string
	to   string
	//conf     *config.Config
}

// NewConsEngine New a consensus engine
func NewConsEngine(from, to string) *ConsEngine {
	ce := new(ConsEngine)
	ce.M = &MsgMapper{MsgMap: make(map[int64]map[string]string)}
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

	if event.Sequence < c.F.sequence {
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

	if n%3 == 0 {
		N = n * 2 / 3
	} else {
		N = n*2/3 + 1
	}

	log.Debugf("[consensus N #%d]", N)
	return int(N)
}

// StartEngine 触发共识引擎尝试处理下一个交易
func (c *ConsEngine) StartEngine() error {

	log.Debugf("Start consensus engine from: [%s] to: [%s] sequence: [%d]", c.from, c.to, c.F.sequence)

	nodes := c.F.conf.GetQscConfig(c.from).NodeAddress

	N := c.consensus32()

	for _, node := range strings.Split(nodes, ",") {

		qcp, err := c.F.queryTxQcpFromNode(c.to, node, c.F.sequence) // be (c.to, node, c.F.sequence)

		if err != nil || qcp == nil {
			continue
		}
		hash := crypto.Sha256(qcp.GetSigData())
		ced := types.CassiniEventDataTx{From: c.from, To: c.to, Height: qcp.BlockHeight, Sequence: c.F.sequence}

		ced.HashBytes = hash

		event := types.Event{NodeAddress: node, CassiniEventDataTx: ced}

		_, err = c.M.AddMsgToMap(c.F, event, N)
		if err != nil {
			return err
		}

	}

	return nil
}

// Ferry Comsumer tx message and handle(consensus, broadcast...) it.
type Ferry struct {
	mtx sync.RWMutex

	sequence int64 //already ferry max sequence

	rmap map[string]*restclient.RestClient //node -> restclient

	conf *config.Config
}

func NewFerry(conf *config.Config, from, to string, sequence int64) *Ferry {

	f := &Ferry{sequence: 0, conf: conf}
	f.rmap = make(map[string]*restclient.RestClient)
	for _, node := range strings.Split(conf.GetQscConfig(from).NodeAddress, ",") {
		add := GetAddressFromUrl(node)
		f.rmap[add] = restclient.NewRestClient(node)

	}
	for _, node := range strings.Split(conf.GetQscConfig(to).NodeAddress, ",") {
		add := GetAddressFromUrl(node)
		f.rmap[add] = restclient.NewRestClient(node)

	}
	return f
}

// SetSequence 设置交易序列号
func (f *Ferry) SetSequence(s int64) {

	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.sequence = s
	log.Infof("sequence set to [#%d]", s)
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

	qscConf := f.conf.GetQscConfig(from)

	// Sign data for public chain
	// Config in QscConfig.Signature
	// true - required
	// false/default - not required
	if qscConf.Signature {
		cdc := app.MakeCodec()
		err = cmn.SignTxQcp(qcp, f.conf.Prikey, cdc)
		if err != nil {
			log.Errorf("Sign Tx Qcp error: %v", err)
		}
		log.Debugf("Sign Tx Qcp for chain: %s", from)
	}

	err = f.postTxQcp(to, qcp)

	if err != nil {
		return errors.New("post qcp transaction failed")
	}

	log.Infof("success ferry qcp transaction from [%s] to [%s] sequence [#%d] \n", from, to, sequence)

	f.SetSequence(f.sequence + 1)

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

	qcp, err = f.queryTxQcpFromNode(to, node, sequence)

	if err != nil || qcp == nil {
		return nil, errors.New("get TxQcp from " + node + "failed.")
	}

	//TODO 取本地联盟链公钥验签
	//pubkey := qcp.Sig.Pubkey  //mock pubkey 为 nil pnic
	//if !pubkey.VerifyBytes(qcp.GetSigData(), qcp.Sig.Signature) {
	//	return nil, errors.New("get TxQcp from " + node + " data verify failed.")
	//}

	// qcp hash 与 hash值比对
	//if string(tmhash.Sum(qcp.GetSigData())) != hash { //算法保持 tmhash.hash 一致 sha256 前 20byte
	hash2 := cmn.Bytes2HexStr(crypto.Sha256(qcp.GetSigData()))
	if hash2 != hash {
		return nil, errors.New("get TxQcp from " + node + "failed")
	}

	return qcp, nil

}

func (f *Ferry) queryTxQcpFromNode(to, node string, sequence int64) (qcp *txs.TxQcp, err error) {

	//"tcp://127.0.0.1:26657"
	//rmap := restclient.NewRestClient(node)
	add := GetAddressFromUrl(node)
	r := f.rmap[add]
	qcp, err = r.GetTxQcp(to, sequence)
	if err != nil && strings.Contains(err.Error(), restclient.ERR_emptyqcp) {
		r := restclient.NewRestClient(node)
		f.rmap[add] = r
		qcp, err = r.GetTxQcp(to, sequence)
	}

	if err != nil || qcp == nil {
		return nil, errors.New("get TxQcp from " + node + "failed.")
	}

	return qcp, nil
}

func GetAddressFromUrl(url string) string {
	n := strings.Index(url, "://")
	if n < 0 {
		return url
	}
	return url[n+3:]
}

func (f *Ferry) postTxQcp(to string, qcp *txs.TxQcp) (err error) {

	success := false
	qscConfig := f.conf.GetQscConfig(to)
	toNodes := qscConfig.NodeAddress
EndPost:
	for _, node := range strings.Split(toNodes, ",") {

		add := GetAddressFromUrl(node)
		r := f.rmap[add]

		err := r.PostTxQcp(to, qcp) //TODO 出错 r := restclient.NewRestClient(node)
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
