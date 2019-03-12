package fabric

import (
	"encoding/json"

	"github.com/QOSGroup/cassini/adapter/ports/fabric/sdk"

	"github.com/QOSGroup/cassini/adapter/ports"
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
	// ChainID id of hyperledger fabric chain
	ChainID = "demo.fabric"
	// ChannelID id of hyperledger fabric chain
	ChannelID = "orgchannel"
	// ChaincodeID id of chaincode
	ChaincodeID = "wallet"
)

type fabricChaincodeQuerySequenceResult struct {
	ChaincodeID string `json:"chaincode,omitempty"`
	Sequence    int64  `json:"sequence,omitempty"`
}

// FabAdaptor provides adapter for hyperledger fabric
type FabAdaptor struct {
	config   *ports.AdapterConfig
	sequence int64
}

// Start fabric adapter service
func (a *FabAdaptor) Start() error {
	return nil
}

// Sync status for fabric adapter service
func (a *FabAdaptor) Sync() error {
	seq, err := a.QuerySequence(a.config.ChainName, "in")
	if err == nil {
		if seq > 1 {
			a.sequence = seq + 1
		} else {
			a.sequence = 1
		}
	}
	return err
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
	return nil
}

// ObtainTx obtain Tx from hyperledger fabric chain
//
// if Tx is register a new account:
//     call ethereum api to create a new account's key, and SubmitTx account's key back to fabric
// if Tx is digital asset withdraw:
//     call ethereum api to transfer
func (a *FabAdaptor) ObtainTx(chainID string, sequence int64) (*txs.TxQcp, error) {
	log.Infof("ObtainTx: %s, %d", chainID, sequence)
	return nil, nil
}

// QuerySequence query sequence of Tx in chaincode
func (a *FabAdaptor) QuerySequence(chainID string, inout string) (int64, error) {
	var as []string
	as = append(as, "sequence", inout)
	args := sdk.Args{
		Func: "query", Args: as}
	var argsArray []sdk.Args
	argsArray = append(argsArray, args)
	ret, err := sdk.ChaincodeQuery(ChannelID, ChaincodeID, argsArray)
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
	log.Infof("query result: %s, %s, %d", chainID, r.ChaincodeID, r.Sequence)
	return r.Sequence, nil
}

// GetSequence returns sequence of tx in cache
func (a *FabAdaptor) GetSequence() int64 {
	return a.sequence
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
