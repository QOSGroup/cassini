package ports

import (
	"errors"
	"fmt"
	"sync"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/route"
	"github.com/QOSGroup/cassini/types"
)

// Builder Create an AdapterService for the specified chain
type Builder func(config AdapterConfig) (AdapterService, error)

// Ports Chain adapter pool interface
type Ports interface {
	Init()
	RegisterBuilder(chain string, builder Builder) error
	Register(config *AdapterConfig) error
	Count(chainName string) int
	Get(chainName string) (map[string]Adapter, error)
}

// defaultPorts Ports default implements
type defaultPorts struct {
	adapters map[string]map[string]AdapterService
	builders map[string]Builder
}

var once sync.Once
var ports Ports

// GetAdapters Get all Adapters for the specified chain
func GetAdapters(chainName string) (map[string]Adapter, error) {
	return GetPortsIncetance().Get(chainName)
}

// RegisterAdapter Check, create and register an Adapter
func RegisterAdapter(config *AdapterConfig) error {
	return GetPortsIncetance().Register(config)
}

// GetPortsIncetance Get Ports singlton instance
func GetPortsIncetance() Ports {
	once.Do(func() {
		ports = &defaultPorts{}
		ports.Init()
	})
	return ports
}

// Init Init the defaultPorts
func (p *defaultPorts) Init() {
	p.adapters = make(map[string]map[string]AdapterService, 0)
	p.builders = make(map[string]Builder, 0)
}

// RegisterBuilder Registers the builder of the adapter for the specified chain
func (p *defaultPorts) RegisterBuilder(chainName string, builder Builder) error {
	if _, ok := p.builders[chainName]; !ok {
		p.builders[chainName] = builder
		return nil
	}
	msg := fmt.Sprintf("builder exist: %s", chainName)
	log.Warnf(msg)
	return errors.New(msg)
}

// Create Check if there is a AdapterService for the specified ip, port and chain-name exist,
// otherwise create one and cache it.
func (p *defaultPorts) Register(conf *AdapterConfig) (err error) {
	log.Infof("Register: %s, %s", conf.ChainName, conf.ChainType)
	var a AdapterService
	adapterKey := GetAdapterKeyByConfig(conf)
	var ads map[string]AdapterService
	var ok bool
	if ads, ok = p.adapters[conf.ChainName]; ok {
		if a, ok = ads[adapterKey]; ok {
			err = fmt.Errorf("adapter already registered: %s", adapterKey)
			return
		}
	}
	if builder, ok := p.builders[conf.ChainType]; ok {
		if conf.Listener == nil {
			nats := config.GetConfig().Queue
			conf.Listener = func(event *types.Event, adapter Adapter) {
				_, err := route.Event2queue(nats, event)
				if err != nil {
					log.Errorf("failed route event to message queue,%s", err.Error())
				}
			}
		}
		if conf.Query == "" {
			conf.Query = "tm.event = 'Tx'  AND qcp.sequence > 0"
		}
		if a, err = builder(*conf); err != nil {
			panic(err)
		}
	} else {
		msg := fmt.Sprintf("no adapter builder found: %s", conf.ChainType)
		log.Warnf(msg)
		return fmt.Errorf(msg)
	}
	if ads, ok = p.adapters[conf.ChainName]; !ok {
		ads = make(map[string]AdapterService, 0)
		p.adapters[conf.ChainName] = ads
	}
	ads[GetAdapterKey(a)] = a
	if err = a.Start(); err != nil {
		panic(err)
	}
	if err = a.Sync(); err != nil {
		panic(err)
	}
	a.Subscribe(conf.Listener)
	return
}

// Count Returns the total number of Adapter for the specified chain-name.
func (p *defaultPorts) Count(chainName string) int {
	return len(p.adapters[chainName])
}

// Apply Apply an unused Adapter for the specified chain-name.
func (p *defaultPorts) Get(chainName string) (map[string]Adapter, error) {
	if adcs, ok := p.adapters[chainName]; ok {
		ads := make(map[string]Adapter, len(adcs))
		for k, v := range adcs {
			ads[k] = v
		}
		return ads, nil
	}
	return nil, errors.New("no adapter found: " + chainName)
}
