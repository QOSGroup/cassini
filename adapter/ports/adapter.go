package ports

import (
	"context"
	"os"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/route"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	tmttypes "github.com/tendermint/tendermint/types"
)

func init() {
	listener := func(event *types.Event, adapter Adapter) {
		_, err := route.Event2queue(config.GetConfig().Nats, event)
		if err != nil {
			log.Errorf("failed route event to message queue,%s", err.Error())
		}
	}
	builder := func(ip string, port int, chain string) (AdapterController, error) {
		a := &qosAdapter{
			chain: chain,
			ip:    ip,
			port:  port}
		a.Start()
		a.Sync()
		a.Subscribe(listener)
		return a, nil
	}
	GetPortsIncetance().RegisterBuilder("qos", builder)
}

// Adapter Chain adapter interface for consensus engine ( consensus.ConsEngine )
// and ferry ( consensus.Ferry )
type Adapter interface {
	SubmitTx(tx *txs.TxQcp) error
	ObtainTx(sequence int64) (*txs.TxQcp, error)
	GetSequence() int64
	// Count Calculate the total and consensus number for chain
	Count() (totalNumber int, consensusNumber int)
	GetChain() string
	GetIP() string
	GetPort() int
}

// EventsListener Listen Tx events from target chain
type EventsListener func(event *types.Event, adapter Adapter)

// AdapterService Chain adapter service interface for adapter pool manager ( adapter.Ports )
type AdapterService interface {
	Start() error
	Sync() error
	Stop() error
	Subscribe(listener EventsListener)
}

/*
AdapterController Inner cache type ( AdapterService and Adapter )

Suitable for a variety of different block chain
*/
type AdapterController interface {
	AdapterService
	Adapter
}

type qosAdapter struct {
	chain    string
	ip       string
	port     int
	sequence int64
	client   *restclient.RestClient
	cancels  []context.CancelFunc
}

func (a *qosAdapter) Start() error {
	a.client = restclient.NewRestClient(GetNodeAddress(a))
	a.cancels = make([]context.CancelFunc, 0)
	return nil
}

func (a *qosAdapter) Sync() error {
	seq, err := a.client.GetSequence(a.chain, "in")
	if err == nil {
		if seq > 1 {
			a.sequence = seq + 1
		} else {
			a.sequence = 1
		}
	}
	return err
}

func (a *qosAdapter) Stop() error {
	if a.client != nil {
		// a.client.close()
	}
	return nil
}

func (a *qosAdapter) Subscribe(listener EventsListener) {
	log.Infof("Starting event subscribe: %s", GetAdapterKey(a))
	remote := "tcp://" + GetNodeAddress(a)
	// go event.EventsSubscribe(remote)
	txs := make(chan interface{})
	go a.subscribeRemote(remote, txs)
	go a.eventHandle(listener, remote, txs)
}

func (a *qosAdapter) SubmitTx(tx *txs.TxQcp) error {
	return nil
}

func (a *qosAdapter) ObtainTx(sequence int64) (qcp *txs.TxQcp, err error) {
	qcp, err = a.client.GetTxQcp(a.chain, sequence)
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

func (a *qosAdapter) GetSequence() int64 {
	return a.sequence
}

func (a *qosAdapter) Count() (totalNumber int, consensusNumber int) {
	totalNumber = GetPortsIncetance().Count(a.chain)
	consensusNumber = Consensus2of3(totalNumber)
	return
}

func (a *qosAdapter) GetChain() string {
	return a.chain
}

func (a *qosAdapter) GetIP() string {
	return a.ip
}

func (a *qosAdapter) GetPort() int {
	return a.port
}

func (a *qosAdapter) subscribeRemote(remote string, txs chan<- interface{}) {
	log.Debug("Event subscribe remote: ", remote)

	//TODO query 条件?? "tm.event = 'Tx' AND qcp.to != '' AND qcp.sequence > 0"
	cancel, err := event.SubscribeRemote(remote,
		"cassini", "tm.event = 'Tx'  AND qcp.sequence > 0", txs)
	if err != nil {
		log.Errorf("Subscibe events failed - remote [%s] : '%s'", remote, err)
		log.Flush()
		os.Exit(1)
	}
	a.cancels = append(a.cancels, cancel)
}

func (a *qosAdapter) eventHandle(listener EventsListener, remote string, txs <-chan interface{}) {
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
