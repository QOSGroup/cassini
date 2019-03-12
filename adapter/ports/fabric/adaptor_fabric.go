package fabric

import (
	"github.com/QOSGroup/cassini/adapter/ports/fabric/sdk"

	"github.com/QOSGroup/cassini/adapter/ports"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
)

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

// Adaptor provides adapter for hyperledger fabric
type Adaptor struct {
	config   *ports.AdapterConfig
	sequence int64
}

// SubmitTx submit Tx to hyperledger fabric chain
func (a *Adaptor) SubmitTx(chainID string, tx *txs.TxQcp) error {
	return nil
}

// ObtainTx obtain Tx from hyperledger fabric chain
//
// if Tx is register a new account:
//     call ethereum api to create a new account's key, and SubmitTx account's key back to fabric
// if Tx is digital asset withdraw:
//     call ethereum api to transfer
func (a *Adaptor) ObtainTx(chainID string, sequence int64) (*txs.TxQcp, error) {
	return nil, nil
}

// QuerySequence query sequence of Tx in chaincode
func (a *Adaptor) QuerySequence(chainID string, inout string) (int64, error) {
	var as []string
	as = append(as, inout)
	args := sdk.Args{
		Func: "query", Args: as}
	var argsArray []sdk.Args
	argsArray = append(argsArray, args)
	ret, err := sdk.ChaincodeQuery(chainID, ChaincodeID, argsArray)
	if err != nil {
		log.Error("query error: %v", err)
		return 0, err
	}
	log.Info("query result: ", ret)
	return 0, nil
}

// GetSequence returns sequence of tx in cache
func (a *Adaptor) GetSequence() int64 {
	return a.sequence
}

// Count Calculate the total and consensus number for chain
func (a *Adaptor) Count() (totalNumber int, consensusNumber int) {
	totalNumber = ports.GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = ports.Consensus2of3(totalNumber)
	return
}

// GetChainName returns chain name
func (a *Adaptor) GetChainName() string {
	return a.config.ChainName
}

// GetIP returns chain node ip
func (a *Adaptor) GetIP() string {
	return a.config.IP
}

// GetPort returns chain node port
func (a *Adaptor) GetPort() int {
	return a.config.Port
}
