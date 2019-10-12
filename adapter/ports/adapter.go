package ports

import (
	"context"
	"fmt"

	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	tctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {
	builder := func(config AdapterConfig) (AdapterService, error) {
		a := &QosAdapter{config: &config}
		return a, nil
	}
	GetPortsIncetance().RegisterBuilder("qos", builder)
}

// Adapter Chain adapter interface for consensus engine ( consensus.ConsEngine )
// and ferry ( consensus.Ferry )
type Adapter interface {
	SubmitTx(chainID string, tx *txs.TxQcp) error
	ObtainTx(chainID string, sequence int64) (*txs.TxQcp, error)
	QuerySequence(chainID string, inout string) (int64, error)
	GetSequence() int64
	// Count Calculate the total and consensus number for chain
	Count() (totalNumber int, consensusNumber int)
	GetChainName() string
	GetIP() string
	GetPort() int
}

// EventsListener Listen Tx events from target chain
type EventsListener func(event *types.Event, adapter Adapter)

// AdapterController Chain adapter controller interface for adapter pool manager ( adapter.Ports )
type AdapterController interface {
	Start() error
	Sync() error
	Stop() error
	Subscribe(listener EventsListener)
}

/*
AdapterService Inner cache type ( AdapterController and Adapter )

Suitable for a variety of different block chain
*/
type AdapterService interface {
	AdapterController
	Adapter
}

// AdapterConfig is parameters for build an AdapterService
type AdapterConfig struct {
	ChainName string
	ChainType string
	IP        string
	Port      int
	Query     string
	Listener  EventsListener
}

// QosAdapter provides adapter for qos chain
type QosAdapter struct {
	config   *AdapterConfig
	sequence int64
	client   *restclient.RestClient
	cancels  []context.CancelFunc
}

// Start qos adapter service
func (a *QosAdapter) Start() error {
	a.client = restclient.NewRestClient(GetNodeAddress(a))
	a.cancels = make([]context.CancelFunc, 0)
	return nil
}

// Sync status for qos adapter service
func (a *QosAdapter) Sync() error {
	seq, err := a.client.GetSequence(a.config.ChainName, "in")
	if err == nil {
		if seq > 1 {
			a.sequence = seq + 1
		} else {
			a.sequence = 1
		}
	}
	return err
}

// Stop qos adapter service
func (a *QosAdapter) Stop() error {
	if a.client != nil {
		// a.client.close()
	}
	return nil
}

// Subscribe events from qos chain
func (a *QosAdapter) Subscribe(listener EventsListener) {
	log.Infof("Starting event subscribe: %s", GetAdapterKey(a))
	remote := "tcp://" + GetNodeAddress(a)
	// go event.EventsSubscribe(remote)
	events := a.subscribeRemote(remote)
	go a.eventHandle(listener, remote, events)
}

// SubmitTx submit Tx to qos chain
func (a *QosAdapter) SubmitTx(chain string, tx *txs.TxQcp) error {
	return a.client.PostTxQcp(chain, tx)
}

// ObtainTx search Tx from qos chain
func (a *QosAdapter) ObtainTx(chain string, sequence int64) (qcp *txs.TxQcp, err error) {
	qcp, err = a.client.GetTxQcp(chain, sequence)
	// if err != nil && !strings.Contains(err.Error(), restclient.ERR_emptyqcp) {
	// 	r := restclient.NewRestClient(node)
	// 	f.rmap[add] = r
	// 	qcp, err = r.GetTxQcp(to, sequence)
	// }
	if err != nil {
		return nil, err
	}

	return qcp, nil
}

// QuerySequence query sequence for the specified chainName and inout ("in" or "out")
func (a *QosAdapter) QuerySequence(chainName string, inout string) (int64, error) {
	return a.client.GetSequence(chainName, inout)
}

// GetSequence returns sequence stored in QosAdapter
func (a *QosAdapter) GetSequence() int64 {
	return a.sequence
}

// Count return total number and consensus number of adapters for qos chain
func (a *QosAdapter) Count() (totalNumber int, consensusNumber int) {
	totalNumber = GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = Consensus2of3(totalNumber)
	return
}

// GetChainName returns chain's name
func (a *QosAdapter) GetChainName() string {
	return a.config.ChainName
}

// GetIP returns chain node's ip
func (a *QosAdapter) GetIP() string {
	return a.config.IP
}

// GetPort returns chain node's port
func (a *QosAdapter) GetPort() int {
	return a.config.Port
}

func (a *QosAdapter) subscribeRemote(remote string) <-chan tctypes.ResultEvent {
	log.Debug("Event subscribe remote: ", remote)
	//TODO query 条件?? "tm.event = 'Tx' AND qcp.to != '' AND qcp.sequence > 0"
	cancel, events, err := event.SubscribeRemote(remote,
		a.config.ChainName, a.config.Query)
	if err != nil {
		// log.Errorf("Subscibe events failed - remote [%s] : '%s'", remote, err)
		// log.Flush()
		// os.Exit(1)
		panic(fmt.Errorf("subscibe events failed: %s", err))
	}
	a.cancels = append(a.cancels, cancel)
	return events
}

func (a *QosAdapter) eventHandle(listener EventsListener, remote string,
	events <-chan tctypes.ResultEvent) {
	for ed := range events {
		log.Debugf("Received event from[%s],'%s'", remote, ed)

		ce := types.CassiniEventDataTx{}
		ce.ConstructFromTags(ed.Events)
		log.Debug("event.Events: ", len(ed.Events), " Height: ", ce.Height)

		event := &types.Event{
			NodeAddress:        remote,
			CassiniEventDataTx: ce,
			Source:             ed}

		listener(event, a)
	}
}
