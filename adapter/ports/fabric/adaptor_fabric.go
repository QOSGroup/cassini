package fabric

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/QOSGroup/cassini/adapter/ports/fabric/sdk"

	"github.com/QOSGroup/cassini/adapter/ports"
	msgtx "github.com/QOSGroup/cassini/adapter/ports/txs"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
)

func init() {
	builder := func(config ports.AdapterConfig) (ports.AdapterService, error) {
		a := &FabAdaptor{config: &config}
		a.Start()
		a.Sync()
		a.Subscribe(config.Listener)
		return a, nil
	}
	ports.GetPortsIncetance().RegisterBuilder("fabric", builder)
}

const (
	// ChainName name of hyperledger fabric
	ChainName = "fabric"
	// ChaincodeID id of chaincode
	ChaincodeID = "wallet"
)

type fabricChaincodeQuerySequenceResult struct {
	ChaincodeID string `json:"chaincode,omitempty"`
	InSequence  int64  `json:"InSequence,omitempty"`
	OutSequence int64  `json:"OutSequence,omitempty"`
}

// ChainResult result of hypelrdger fabric chaincode
type ChainResult struct {
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	ErrString string `json:"error,omitempty"`
	// Result    interface{} `json:"result,omitempty"`
}

// FabAdaptor provides adapter for hyperledger fabric
type FabAdaptor struct {
	config      *ports.AdapterConfig
	inSequence  int64
	outSequence int64
}

// Start fabric adapter service
func (a *FabAdaptor) Start() error {
	return nil
}

// Sync status for fabric adapter service
func (a *FabAdaptor) Sync() error {
	// seq, err := a.QuerySequence(a.config.ChainName, "in")
	// if err == nil {
	// 	if seq > 1 {
	// 		a.outSequence = seq + 1
	// 	} else {
	// 		a.outSequence = 1
	// 	}
	// }
	return nil
}

// Stop fabric adapter service
func (a *FabAdaptor) Stop() error {
	return nil
}

// Subscribe events from fabric chain
func (a *FabAdaptor) Subscribe(listener ports.EventsListener) {
	log.Infof("no event subscribe: %s", ports.GetAdapterKey(a))
}

// SubmitTx submit Tx to hyperledger fabric chain
func (a *FabAdaptor) SubmitTx(chainID string, tx *txs.TxQcp) error {
	bytes := tx.TxStd.ITx.GetSignData()
	log.Infof("SubmitTx: %s(%s) %d register block: %s",
		a.GetChainName(), chainID, tx.Sequence, string(bytes))
	var args []string
	args = append(args, "block", string(bytes))
	arg := sdk.Args{Func: "register", Args: args}
	var argsArray []sdk.Args
	argsArray = append(argsArray, arg)
	ret, err := sdk.ChaincodeInvoke(ChaincodeID, argsArray)
	if err != nil {
		log.Errorf("SubmitTx: %s(%s) %d: register block error: %v",
			a.GetChainName(), chainID, tx.Sequence, err)
		return err
	}
	log.Infof("SubmitTx: %s(%s) %d: register block response: %s",
		a.GetChainName(), chainID, tx.Sequence, ret)
	result := &ChainResult{}
	if err = json.Unmarshal([]byte(ret), result); err != nil {
		log.Errorf("SubmitTx: %s(%s) %d: register block result unmarshal error: %v",
			a.GetChainName(), chainID, tx.Sequence, err)
		return err
	}
	if result.Code != http.StatusOK {
		log.Errorf("SubmitTx: %s(%s) %d: register block failed: %s",
			a.GetChainName(), chainID, tx.Sequence, ret)
		err = fmt.Errorf("SubmitTx: %s(%s) %d: register block failed: %s",
			a.GetChainName(), chainID, tx.Sequence, ret)
		return err
	}
	return nil
}

// ObtainTx obtain Tx from hyperledger fabric chain
//
// if Tx is register a new account:
//     call ethereum api to create a new account's key, and SubmitTx account's key back to fabric
// if Tx is digital asset withdraw:
//     call ethereum api to transfer
func (a *FabAdaptor) ObtainTx(chainID string, sequence int64) (*txs.TxQcp, error) {
	log.Infof("ObtainTx: %s(%s), %d", a.GetChainName(), chainID, sequence)
	var as []string
	as = append(as, "transaction", "sequence", strconv.FormatInt(sequence, 10))
	args := sdk.Args{
		Func: "query", Args: as}
	var argsArray []sdk.Args
	argsArray = append(argsArray, args)
	ret, err := sdk.ChaincodeQuery(ChaincodeID, argsArray)
	if err != nil {
		log.Errorf("ObtainTx %s(%s), %d error: %v",
			a.GetChainName(), chainID, sequence, err)
		return nil, err
	}
	log.Info("query transaction result: ", ret)
	tx := msgtx.NewTxQcp(fmt.Sprintf("%s(%s)", a.GetChainName(), chainID),
		a.GetChainName(), chainID, int64(1), int64(sequence), ret)
	return tx, nil
}

// QuerySequence query sequence of Tx in chaincode
func (a *FabAdaptor) QuerySequence(chainID string, inout string) (int64, error) {
	var as []string
	as = append(as, "sequence")
	args := sdk.Args{
		Func: "query", Args: as}
	var argsArray []sdk.Args
	argsArray = append(argsArray, args)
	ret, err := sdk.ChaincodeQuery(ChaincodeID, argsArray)
	if err != nil {
		log.Error("query error: %v", err)
		return 0, err
	}
	r := &fabricChaincodeQuerySequenceResult{}
	err = json.Unmarshal([]byte(ret), r)
	if err != nil {
		log.Errorf("parse chain result error: %v\n%s", err, ret)
		return 0, err
	}
	a.inSequence = r.InSequence
	a.outSequence = r.OutSequence
	if strings.EqualFold("in", inout) {
		log.Infof("QuerySequence: %s(%s), in %d",
			a.GetChainName(), chainID, r.InSequence)
		return r.InSequence, nil
	}
	log.Infof("QuerySequence: %s(%s), out %d",
		a.GetChainName(), chainID, r.OutSequence)
	return r.OutSequence, nil
}

// GetSequence returns sequence of tx in cache
func (a *FabAdaptor) GetSequence() int64 {
	return a.outSequence
}

// Count Calculate the total and consensus number for chain
func (a *FabAdaptor) Count() (totalNumber int, consensusNumber int) {
	totalNumber = ports.GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = ports.Consensus2of3(totalNumber)
	return
}

// GetChainName returns chain name
func (a *FabAdaptor) GetChainName() string {
	return a.config.ChainName
}

// GetIP returns chain node ip
func (a *FabAdaptor) GetIP() string {
	return a.config.IP
}

// GetPort returns chain node port
func (a *FabAdaptor) GetPort() int {
	return a.config.Port
}
