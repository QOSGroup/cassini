package ethereum

import (
	"encoding/json"
	"strings"

	"github.com/QOSGroup/cassini/adapter/ports"
	fabricTx "github.com/QOSGroup/cassini/adapter/ports/fabric/sdk/tx"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
)

func init() {
	builder := func(config ports.AdapterConfig) (ports.AdapterService, error) {
		a := &EthAdaptor{config: &config}
		a.Start()
		a.Sync()
		a.Subscribe(config.Listener)
		return a, nil
	}
	ports.GetPortsIncetance().RegisterBuilder("ethereum", builder)
}

// EthAdaptor provides adapter for ethereum
type EthAdaptor struct {
	config      *ports.AdapterConfig
	inSequence  int64
	outSequence int64
}

// Start fabric adapter service
func (a *EthAdaptor) Start() error {
	return nil
}

// Sync status for fabric adapter service
func (a *EthAdaptor) Sync() error {
	seq, err := a.QuerySequence(a.config.ChainName, "in")
	if err == nil {
		if seq > 1 {
			a.outSequence = seq + 1
		} else {
			a.outSequence = 1
		}
	}
	return err
}

// Stop fabric adapter service
func (a *EthAdaptor) Stop() error {
	return nil
}

// Subscribe events from ethereum chain
func (a *EthAdaptor) Subscribe(listener ports.EventsListener) {
	log.Infof("no event subscribe: %s", ports.GetAdapterKey(a))
}

// SubmitTx submit Tx to ethereum chain
func (a *EthAdaptor) SubmitTx(chainID string, tx *txs.TxQcp) error {
	jsonTx := tx.TxStd.ITx.GetSignData()
	log.Infof("SubmitTx: %s(%s) %d: chain result: %s",
		a.GetChainName(), chainID, tx.Sequence, jsonTx)
	t := fabricTx.WalletTx{}
	err := json.Unmarshal([]byte(jsonTx), &t)
	if err != nil {
		log.Errorf("SubmitTx: %s(%s) error: %v",
			a.GetChainName(), chainID, err)
		return err
	}
	a.inSequence = tx.Sequence
	if a.outSequence <= 1 {
		a.outSequence = t.Height
	}
	// encrypted
	// etcd
	// (recharge) query ethereum transactions
	// (withdraw) transfer
	return nil
}

// ObtainTx obtain Tx from ethereum chain
// recharge:
//     send transaction data back to fabric
// withdraw:
//     send transaction data back to fabric
func (a *EthAdaptor) ObtainTx(chainID string, sequence int64) (*txs.TxQcp, error) {
	log.Infof("ObtainTx: %s(%s), %d", a.GetChainName(), chainID, sequence)
	return nil, nil
}

// QuerySequence query sequence of Tx in ethereum
func (a *EthAdaptor) QuerySequence(chainID string, inout string) (int64, error) {
	if strings.EqualFold("in", inout) {
		log.Infof("QuerySequence: %s(%s), in %d",
			a.GetChainName(), chainID, a.inSequence)
		return a.inSequence, nil
	}
	log.Infof("QuerySequence: %s(%s), out %d",
		a.GetChainName(), chainID, a.outSequence)
	return a.outSequence, nil
}

// GetSequence returns sequence of tx in cache
func (a *EthAdaptor) GetSequence() int64 {
	return a.outSequence
}

// Count Calculate the total and consensus number for chain
func (a *EthAdaptor) Count() (totalNumber int, consensusNumber int) {
	totalNumber = ports.GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = ports.Consensus2of3(totalNumber)
	log.Debugf("total: %d; consensus: %d;", totalNumber, consensusNumber)
	return
}

// GetChainName returns chain name
func (a *EthAdaptor) GetChainName() string {
	return a.config.ChainName
}

// GetIP returns chain node ip
func (a *EthAdaptor) GetIP() string {
	return a.config.IP
}

// GetPort returns chain node port
func (a *EthAdaptor) GetPort() int {
	return a.config.Port
}
