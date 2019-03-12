package sdk

import (
	"github.com/ethereum/go-ethereum/dashboard"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
)

const (
	// VERSION of ethereum
	VERSION = "v1.8.23"
	// CLIENT identifier
	CLIENT = "geth"
)

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = CLIENT
	cfg.Version = VERSION
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "shh")
	cfg.WSModules = append(cfg.WSModules, "eth", "shh")
	cfg.IPCPath = "geth.ipc"
	return cfg
}

// GethConfig wrapped ethereum config
type GethConfig struct {
	Eth  eth.Config
	Shh  whisper.Config
	Node node.Config
	// Ethstats  ethstatsConfig
	Dashboard dashboard.Config
}

// Config returns ethereum config
func Config() GethConfig {
	return GethConfig{Node: defaultNodeConfig()}
}
