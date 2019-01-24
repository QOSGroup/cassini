package ports

import (
	"context"
	"os"

	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	tmttypes "github.com/tendermint/tendermint/types"
)

// QosBuilder is default builder for qos chain
var QosBuilder Builder = func(config AdapterConfig) (AdapterService, error) {
	a := &QosAdapter{config: &config}
	a.Start()
	a.Sync()
	a.Subscribe(config.Listener)
	return a, nil
}

func init() {
	GetPortsIncetance().RegisterBuilder("qos", QosBuilder)
}

// Adapter Chain adapter interface for consensus engine ( consensus.ConsEngine )
// and ferry ( consensus.Ferry )
type Adapter interface {
	SubmitTx(tx *txs.TxQcp) error
	ObtainTx(sequence int64) (*txs.TxQcp, error)
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
type AdapterConfig struct{
	ChainName string
	IP       string
	Port     int
	Query string
	Listener EventsListener
}

type QosAdapter struct {
	config *AdapterConfig
	sequence int64
	client   *restclient.RestClient
	cancels  []context.CancelFunc
}

func (a *QosAdapter) Start() error {
	a.client = restclient.NewRestClient(GetNodeAddress(a))
	a.cancels = make([]context.CancelFunc, 0)
	return nil
}

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

func (a *QosAdapter) Stop() error {
	if a.client != nil {
		// a.client.close()
	}
	return nil
}

func (a *QosAdapter) Subscribe(listener EventsListener) {
	log.Infof("Starting event subscribe: %s", GetAdapterKey(a))
	remote := "tcp://" + GetNodeAddress(a)
	// go event.EventsSubscribe(remote)
	txs := make(chan interface{})
	go a.subscribeRemote(remote, txs)
	go a.eventHandle(listener, remote, txs)
}

func (a *QosAdapter) SubmitTx(tx *txs.TxQcp) error {
	return nil
}

func (a *QosAdapter) ObtainTx(sequence int64) (qcp *txs.TxQcp, err error) {
	qcp, err = a.client.GetTxQcp(a.GetChainName(), sequence)
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

func (a *QosAdapter) GetSequence() int64 {
	return a.sequence
}

func (a *QosAdapter) Count() (totalNumber int, consensusNumber int) {
	totalNumber = GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = Consensus2of3(totalNumber)
	return
}

func (a *QosAdapter) GetChainName() string {
	return a.config.ChainName
}

func (a *QosAdapter) GetIP() string {
	return a.config.IP
}

func (a *QosAdapter) GetPort() int {
	return a.config.Port
}

func (a *QosAdapter) subscribeRemote(remote string, txs chan<- interface{}) {
	log.Debug("Event subscribe remote: ", remote)
	//TODO query 条件?? "tm.event = 'Tx' AND qcp.to != '' AND qcp.sequence > 0"
	cancel, err := event.SubscribeRemote(remote,
		a.config.ChainName, a.config.Query, txs)
	if err != nil {
		log.Errorf("Subscibe events failed - remote [%s] : '%s'", remote, err)
		log.Flush()
		os.Exit(1)
	}
	a.cancels = append(a.cancels, cancel)
}

func (a *QosAdapter) eventHandle(listener EventsListener, remote string, txs <-chan interface{}) {
	for ed := range txs {
		edt := ed.(tmttypes.EventDataTx)
		log.Debugf("Received event from[%s],'%s'", remote, edt)

		cassiniEventDataTx := types.CassiniEventDataTx{}
		cassiniEventDataTx.Height = edt.Height
		cassiniEventDataTx.ConstructFromTags(edt.Result.Tags)

		event := types.Event{
			NodeAddress:        remote,
			CassiniEventDataTx: cassiniEventDataTx}

		listener(&event, a)
	}
}
