package ports

import (
	"errors"
	"fmt"
	"sync"

	"github.com/QOSGroup/cassini/log"
)

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

// Builder Create an Adapter
type Builder func(ip string, port int, chain string) (AdapterController, error)

// Init Init the defaultPorts
func (p *defaultPorts) Init() {
	p.adapters = make(map[string]map[string]AdapterController, 0)
	p.builders = make(map[string]Builder, 0)
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
