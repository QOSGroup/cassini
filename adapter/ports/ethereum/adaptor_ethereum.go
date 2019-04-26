package ethereum

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/QOSGroup/cassini/adapter/ports"
	"github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
	fabsdk "github.com/QOSGroup/cassini/adapter/ports/fabric/sdk"
	fabricTx "github.com/QOSGroup/cassini/adapter/ports/fabric/sdk/tx"
	msgtx "github.com/QOSGroup/cassini/adapter/ports/txs"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/storage"
	"github.com/QOSGroup/qbase/txs"
)

func init() {
	builder := func(config ports.AdapterConfig) (ports.AdapterService, error) {
		a := &EthAdaptor{
			config: &config,
			Addrs:  storage.NewAddressCacheBook()}
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
	Addrs       storage.AddressBook
}

// Start fabric adapter service
func (a *EthAdaptor) Start() error {
	return nil
}

// Sync status for fabric adapter service
func (a *EthAdaptor) Sync() error {
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
		height, err := strconv.ParseInt(t.Height[2:], 16, 64)
		if err != nil {
			log.Errorf("SubmitTx: %s(%s) parse height(%s) error: %v",
				a.GetChainName(), chainID, t.Height, err)
			return err
		}
		a.outSequence = height
	}
	// cache wallet address for block filtering in ObtainTx
	if strings.EqualFold("register", t.Func) {
		a.Addrs.Add(t.Address)
		log.Infof("new ethereum wallet address: %s", t.Address)
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
	// ignore useless block, 10000000 in kovan
	if sequence < 10000000 {
		return nil, errors.New("invalid sequence")
	}
	block, err := sdk.EthGetBlockByNumber(sequence)
	if err != nil {
		log.Errorf("ethereum rpc error: %v", err)
		return nil, err
	}
	registerBlock := &fabsdk.BlockRegister{
		Height: block.Number}
	for _, tx := range block.Transactions {
		if a.Addrs.Exist(tx.To) ||
			a.Addrs.Exist(tx.From) {
			receipt, err := sdk.EthGetTransactionReceipt(tx.Hash)
			if err != nil {
				log.Errorf("check block transaction receipt hash: %s error: %v",
					tx.Hash, err)
				return nil, err
			}
			t := &fabsdk.TxRegister{
				Chain:    "ethereum",
				Token:    "eth",
				From:     tx.From,
				To:       tx.To,
				Txhash:   tx.Hash,
				GasUsed:  receipt.GasUsed,
				GasPrice: tx.GasPrice,
				Status:   receipt.Status}
			if receipt.Success() {
				t.Amount = tx.Value
			} else {
				log.Warnf("ObtainTx: %s(%s) %d check block: %s",
					a.GetChainName(), chainID, sequence,
					fmt.Sprintf("transaction reverted hash: %s", tx.Hash))
			}
			registerBlock.Txs = append(registerBlock.Txs, t)
		}
	}
	bytes, err := json.Marshal(registerBlock)
	if err != nil {
		log.Errorf("check block marshal error: %v", err)
		return nil, err
	}
	jsonRegisterBlock := string(bytes)
	log.Infof("ObtainTx: %s(%s) %d check block: %s", a.GetChainName(), chainID,
		sequence, jsonRegisterBlock)
	tx := msgtx.NewTxQcp(fmt.Sprintf("%s(%s)", a.GetChainName(), chainID),
		a.GetChainName(), chainID, int64(1), int64(sequence), jsonRegisterBlock)
	// a.outSequence = sequence + 1
	return tx, nil
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
	log.Debugf("%s adaptor count: %d; consensus: %d;",
		a.GetChainName(), totalNumber, consensusNumber)
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
