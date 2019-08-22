package consensus

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/QOSGroup/qbase/example/basecoin/app"
	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/tendermint/crypto"

	"errors"
	"fmt"
	"time"

	"github.com/QOSGroup/cassini/adapter/ports"
	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/concurrency"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	exporter "github.com/QOSGroup/cassini/prometheus"
	"github.com/QOSGroup/cassini/restclient"
)

// Ferry Comsumer tx message and handle(consensus, broadcast...) it.
type Ferry struct {
	mtx sync.RWMutex

	mutex    concurrency.Mutex
	from, to string
	sequence int64 //already ferry max sequence

	ConsMap *ConsensusMap
	conf    *config.Config
}

// NewFerry new a ferry for relaying transactions
func NewFerry(conf *config.Config, from, to string, sequence int64) *Ferry {

	f := &Ferry{sequence: 1, conf: conf}
	f.from, f.to = from, to
	f.ConsMap = &ConsensusMap{ConsMap: make(map[int64]map[string]string)}

	seq, _ := f.GetSequenceFromChain(from, to, "in")
	if seq > 1 {
		f.sequence = seq + 1
	} else {
		f.sequence = 1
	}
	if f.conf.UseEtcd {
		var err error
		f.mutex, err = concurrency.NewMutex(from+"_"+to, f.conf.Lock)
		if err != nil {
			log.Errorf("create mutex failed. %v", err)
			f.conf.UseEtcd = false
			return f
		}
		f.mutex.Update(f.sequence)
	}

	return f
}

// StartFerry 启动共识引擎尝试处理下一个交易
func (f *Ferry) StartFerry() error {
	go func() {
		for {
			f.SetSequence(f.from, f.to, 0)
			seq := atomic.LoadInt64(&f.sequence)
			cons, err := f.ConsMap.GetConsFromMap(seq)

			if err != nil {
				log.Warnf("consensus sequence: [#%d] error: %v", seq, err)
				time.Sleep(time.Duration(f.conf.EventWaitMillitime) * time.Millisecond)
			} else if cons != nil { //已有该sequence 共识
				if err = f.ferryQCP(cons.Hash, cons.Nodes, seq); err != nil {
					log.Errorf("ferry qcp transaction f.t.s[%s %s #%d] hash[%s] failed. %v",
						f.from, f.to, seq, cons.Hash[:10], err)
					time.Sleep(time.Duration(f.conf.EventWaitMillitime) * time.Millisecond)
				}
			}

		}
	}()
	return nil
}

// SetSequence 设置交易序列号
func (f *Ferry) SetSequence(from, to string, s int64) {

	f.mtx.Lock()
	defer f.mtx.Unlock()

	seq, _ := f.GetSequenceFromChain(from, to, "in")
	log.Infof("f.t[%s %s] ferry set sequence: %d or %d?", from, to, s, seq)
	// f.sequence = common.MaxInt64(s, seq) + 1
	// TODO ??? why max + 1?
	f.sequence = seq + 1

	log.Infof("f.t[%s %s] ferry sequence set to [#%d]", from, to, f.sequence)
}

//GetSequenceFromChain 在to chain上查询 来自/要去 from chain 的 sequence
func (f *Ferry) GetSequenceFromChain(from, to, inout string) (int64, error) {

	// qsc := f.conf.GetQscConfig(to)

	// nodeto := strings.Split(qsc.NodeAddress, ",")

	// var seq int64
	// var err error
	// for _, n := range nodeto {
	// 	add := GetAddressFromUrl(n)
	// 	r := f.rmap[add]
	// 	seq, err = r.GetSequence(from, inout)
	// 	if err == nil {
	// 		return seq, nil
	// 	}
	// }
	var seq int64
	ads, err := ports.GetAdapters(to)
	if err != nil {
		return -1, err
	}
	for _, a := range ads {
		seq, err = a.QuerySequence(from, inout)
		if err == nil {
			return seq, nil
		}
	}
	return -1, err
}

//ferryQCP get qcp transaction from source chain and post it to destnation chain
//
//from is chain name of the source chain
//to is the chain name of destnation chain
//nodes is consensus nodes of the source chain
func (f *Ferry) ferryQCP(hash, nodes string, sequence int64) (err error) {

	log.Debugf("ferry qcp transaction f.t.s[%s %s %d]", f.from, f.to, sequence)

	qcp, err := f.getTxQcp(f.from, f.to, hash, nodes, sequence)

	if err != nil { //TODO 拜占庭共识失败后 循环至此

		//log.Errorf("ferry qcp transaction from [%s] to [%s] sequence [%d]. %s", f.from, f.to, sequence, err.Error())
		return errors.New("get qcp transaction failed," + err.Error())
	}

	qscConf := f.conf.GetQscConfig(f.from)

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
		log.Debugf("Sign Tx Qcp for chain: %s", f.from)
	}

	if f.conf.UseEtcd {
		if lockseq, err := f.mutex.Lock(sequence); err != nil {
			log.Warnf("get lock failed. sequence [%d] / [%d]", sequence, lockseq)
			if lockseq > 0 {
				if sequence < lockseq {
					log.Warnf("update sequence [#%d] to etcd sequence [#%d]",
						sequence, lockseq)
					// f.SetSequence(f.from, f.to, lockseq-1)
				}
			}
			return fmt.Errorf("get lock fail %v", err)
		}
		log.Debugf("get lock success ,sequence [%d]", sequence)
		err = f.postTxQcp(f.to, qcp)

		if err != nil {
			f.mutex.Unlock(false)
			//log.Errorf("post qcp transaction failed. %v", err)
			return errors.New("post qcp transaction failed")
		}
		f.mutex.Unlock(true)
	} else {
		err = f.postTxQcp(f.to, qcp)
		if err != nil {
			log.Errorf("post qcp transaction failed. %v", err)
			return errors.New("post qcp transaction failed")
		}
	}
	delete(f.ConsMap.ConsMap, sequence)
	log.Infof("success ferry qcp transaction f.t.s[%s %s #%d] \n",
		f.from, f.to, sequence)

	// TODO check sequence at everytime???
	// f.SetSequence(f.from, f.to, f.sequence)

	return nil

}

//getTxQcp get QCP transactions from sorce chain
func (f *Ferry) getTxQcp(from, to, hash, nodes string, sequence int64) (qcp *txs.TxQcp, err error) {

	success := false
	bempty := false

EndGet:

	for _, node := range strings.Split(nodes, ",") {

		qcp, err = f.getTxQcpFromNode(from, to, hash, node, sequence)

		if err != nil || qcp == nil {
			if !strings.Contains(err.Error(), restclient.ERR_emptyqcp) {
				log.Warnf("get transaction from %s failed,%s", node, err.Error())
			} else {
				bempty = true
			}

			continue
		}

		success = true
		break EndGet

	}

	if !success {
		if bempty {
			return nil, errors.New(restclient.ERR_emptyqcp)
		}
		return nil, errors.New("getTxQcp failed")
	}

	return qcp, nil
}

// func (f *Ferry) getTxQcpParalle(from, to, hash, nodes string, sequence int64) (qcps []txs.TxQcp, err error) {

// 	nodeList := strings.Split(nodes, ",")
// 	var tasks = make([]common.Task, len(nodeList))

// 	for i := 0; i < len(tasks); i++ {
// 		tasks[i] = func(i int) (res interface{}, err error, abort bool) {
// 			qcp, err := f.getTxQcpFromNode(to, hash, nodeList[i], sequence)
// 			return qcp, err, false //TODO
// 		}
// 	}

// 	var tResults, ok = common.Parallel(tasks...)
// 	if !ok {
// 		log.Error("parallel failed")
// 	}

// 	var failTasks int
// 	for i := 0; i < len(tasks); i++ {
// 		tResult, ok := tResults.LatestResult(i)
// 		if !ok {
// 			failTasks++
// 		} else if tResult.Error != nil {
// 			failTasks++
// 		} else {
// 			qcps = append(qcps, *(tResult.Value).(*txs.TxQcp))
// 		}

// 	}

// 	if len(qcps)*2 > failTasks { //TODO 加入共识逻辑
// 		return qcps, nil
// 	}

// 	return nil, errors.New("parallel get qcp transaction from chain " + from + " failed")
// }

//getTxQcpFromNode get QCP transactions from single chain node
func (f *Ferry) getTxQcpFromNode(from, to, hash, node string, sequence int64) (qcp *txs.TxQcp, err error) {

	qcp, err = f.queryTxQcpFromNode(from, to, node, sequence)

	if err != nil || qcp == nil {
		return nil, err
	}

	//TODO 取本地联盟链公钥验签
	//pubkey := qcp.Sig.Pubkey  //mock pubkey 为 nil pnic
	//if !pubkey.VerifyBytes(qcp.GetSigData(), qcp.Sig.Signature) {
	//	return nil, errors.New("get TxQcp from " + node + " data verify failed.")
	//}

	// qcp hash 与 hash值比对
	//if string(tmhash.Sum(qcp.GetSigData())) != hash { //算法保持 tmhash.hash 一致 sha256 前 20byte
	hash2 := cmn.Bytes2HexStr(crypto.Sha256(qcp.BuildSignatureBytes()))
	if hash2 != hash {
		return nil, fmt.Errorf("get TxQcp failed, tansaction hash not correct: %s",
			fmt.Sprintf("node: %s from: %s to: %s sequence: %d hash: %s hash2: %s",
				node, from, to, sequence, hash, hash2))
	}

	return qcp, nil

}

//queryTxQcpFromNode get TxQcp from node
//
//to destnation chain id
func (f *Ferry) queryTxQcpFromNode(from, to, node string, sequence int64) (qcp *txs.TxQcp, err error) {

	// add := GetAddressFromUrl(node)
	// r := f.rmap[add]
	// qcp, err = r.GetTxQcp(to, sequence)
	// if err != nil && !strings.Contains(err.Error(), restclient.ERR_emptyqcp) {
	// 	r := restclient.NewRestClient(node)
	// 	f.rmap[add] = r
	// 	qcp, err = r.GetTxQcp(to, sequence)
	// }

	// if err != nil || qcp == nil {
	// 	return nil, err
	// }

	// return qcp, nil
	var ads map[string]ports.Adapter
	ads, err = ports.GetAdapters(from)
	log.Warnf("queryTxQcpFromNode: from %s: to %s: adapters: %d: seq: %d", from, to, len(ads), sequence)
	if err == nil {
		for _, a := range ads {
			qcp, err = a.ObtainTx(to, sequence)
			if err == nil {
				return
			}
		}
	}
	return
}

func (f *Ferry) postTxQcp(to string, qcp *txs.TxQcp) (err error) {
	ads, err := ports.GetAdapters(to)
	if err == nil {
		for _, a := range ads {
			log.Debug("post TxQcp to: ", to)
			err := a.SubmitTx(to, qcp)
			if err != nil {
				log.Errorf("post TxQcp error: %v", err)
			} else {
				exporter.Set(exporter.KeyTxsPerSecond, 1)
				return nil
			}
		}
		err = errors.New("post qcp transaction failed")
	}
	return err
}
