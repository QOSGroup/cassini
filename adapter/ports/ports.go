package ports

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/route"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	tmttypes "github.com/tendermint/tendermint/types"
)

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
	Subcribe(listener EventsListener)
}

/*
AdapterController Inner cache type ( AdapterService and Adapter )

Suitable for a variety of different block chain
*/
type AdapterController interface {
	AdapterService
	Adapter
}

// Ports Chain adapter pool interface
type Ports interface {
	Init()
	RegisterBuilder(chain string, builder Builder) error
	Register(ip string, port int, chain string) error
	Count(chain string) int
	Get(chain string) (map[string]Adapter, error)
}

// defaultPorts Ports default implements
type defaultPorts struct {
	adapters map[string]map[string]AdapterController
	builders map[string]Builder
}

var once sync.Once
var ports Ports

// GetAdapters Get all Adapters for the specified chain
func GetAdapters(chain string) (map[string]Adapter, error) {
	return GetPortsIncetance().Get(chain)
}

// RegisterAdapter Check, create and register an Adapter
func RegisterAdapter(ip string, port int, chain string) error {
	return GetPortsIncetance().Register(ip, port, chain)
}

// GetPortsIncetance Get Ports singlton instance
func GetPortsIncetance() Ports {
	once.Do(func() {
		ports = &defaultPorts{}
		ports.Init()
	})
	return ports
}

// GetAdapterKey Gen an adapter key
func GetAdapterKey(a Adapter) string {
	return fmt.Sprintf("%s://%s:%d", a.GetChain(), a.GetIP(), a.GetPort())
}

// GetNodeAddress Gen a node address
func GetNodeAddress(a Adapter) string {
	return fmt.Sprintf("%s:%d", a.GetIP(), a.GetPort())
}

// Consensus2of3 Calculate number of consensus
func Consensus2of3(value int) int {
	return (value*2 + 2) / 3
}

// Builder Create an Adapter
type Builder func(ip string, port int, chain string) (AdapterController, error)

// Init Init the defaultPorts
func (p *defaultPorts) Init() {
	p.adapters = make(map[string]map[string]AdapterController, 0)
	p.builders = make(map[string]Builder, 0)

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
		a.Subcribe(listener)
		return a, nil
	}
	p.RegisterBuilder("qos", builder)

}

// RegisterBuilder Registers the builder of the adapter for the specified chain
func (p *defaultPorts) RegisterBuilder(chain string, builder Builder) error {
	if _, ok := p.builders[chain]; !ok {
		p.builders[chain] = builder
		return nil
	}
	log.Warnf("builder exist: %s", chain)
	return fmt.Errorf("builder exist: %s", chain)
}

// Create Check if there is a Adapter for the specified ip, port and chain exist,
// otherwise create one and cache it.
func (p *defaultPorts) Register(ip string, port int, chain string) (err error) {
	var a AdapterController
	chain = "qos"
	if builder, ok := p.builders[chain]; ok {
		a, err = builder(ip, port, chain)
	} else {
		log.Warnf("no adapter builder found: %s", chain)
		return fmt.Errorf("no adapter builder found: %s", chain)
	}
	var ads map[string]AdapterController
	var ok bool
	if ads, ok = p.adapters[chain]; !ok {
		ads = make(map[string]AdapterController, 0)
		p.adapters[chain] = ads
	}
	ads[GetAdapterKey(a)] = a
	return nil
}

// Count Returns the total number of Adapter for the specified chain.
func (p *defaultPorts) Count(chain string) int {
	return len(p.adapters[chain])
}

// Apply Apply an unused Adapter for the specified chain.
func (p *defaultPorts) Get(chain string) (map[string]Adapter, error) {
	if adcs, ok := p.adapters[chain]; ok {
		ads := make(map[string]Adapter, len(adcs))
		for k, v := range adcs {
			ads[k] = v
		}
		return ads, nil
	}
	return nil, errors.New("no adapter found: " + chain)
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

func (a *qosAdapter) Subcribe(listener EventsListener) {
	log.Infof("Starting event subcribe: %s", GetAdapterKey(a))
	remote := "tcp://" + GetNodeAddress(a)
	// go event.EventsSubscribe(remote)
	txs := make(chan interface{})
	go func() {
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
	}()
	go func() {
		for ed := range txs {

			edt := ed.(tmttypes.EventDataTx)
			log.Debugf("Received event from[%s],'%s'", remote, edt)

			cassiniEventDataTx := types.CassiniEventDataTx{}
			cassiniEventDataTx.Height = edt.Height
			cassiniEventDataTx.ConstructFromTags(edt.Result.Tags)

			event := types.Event{NodeAddress: remote, CassiniEventDataTx: cassiniEventDataTx}

			listener(&event, a)
		}
	}()
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
