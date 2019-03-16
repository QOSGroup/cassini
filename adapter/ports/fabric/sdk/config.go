package sdk

import (
	"strings"
	"sync"
	"time"

	"github.com/QOSGroup/cassini/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	sdkconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	secAction "github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/action"
	cliconfig "github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/config"
	"github.com/spf13/pflag"
)

const (
	// AutoDetectSelectionProvider indicates that a selection provider is to be automatically determined using channel capabilities
	AutoDetectSelectionProvider = "auto"
)

// FabConfig wrap configuration values
type FabConfig struct {
	ConfigFile         string
	UserName           string
	ChainID            string
	ChannelID          string
	SelectionProvider  string
	OrdererURL         string
	IsLoggingEnabledFo bool                // enable log
	ConfigProvider     core.ConfigProvider // provider returns the config provider
	Concurrency        uint16
	OrgIDs             []string
	PeerURL            string
	PeerURLs           []string
	MaxAttempts        int
	InitialBackoff     time.Duration
	MaxBackoff         time.Duration
	BackoffFactor      float64
	Verbose            bool
	Iterations         int
	PrintPayloadOnly   bool
	PrintFormat        string
	Writer             string
	Base64             bool
}

var once *sync.Once
var config *FabConfig

var peerURL = "localhost:7051,localhost:8051"
var orgIDsStr = "localhost:7050"

func init() {
	once = &sync.Once{}
}

// Config get config's singleton
func Config() *FabConfig {
	once.Do(func() {
		config = &FabConfig{
			ConfigFile:         "/vagrant/gopath/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml",
			UserName:           "",
			ChainID:            "demo.fabric",
			ChannelID:          "orgchannel",
			SelectionProvider:  AutoDetectSelectionProvider,
			OrdererURL:         "",
			IsLoggingEnabledFo: true,
			Concurrency:        1,
			PeerURL:            "localhost:7051",
			MaxAttempts:        3,
			InitialBackoff:     time.Duration(1000) * time.Millisecond,
			MaxBackoff:         time.Duration(5000) * time.Millisecond,
			BackoffFactor:      2,
			Verbose:            false,
			Iterations:         1,
			PrintPayloadOnly:   false,
			PrintFormat:        "json",
			Writer:             "stdout",
			Base64:             false,
		}
		urls := parse(peerURL, ",")
		config.PeerURLs = urls
		urls = parse(orgIDsStr, ",")
		config.OrgIDs = urls
		config.ConfigProvider = sdkconfig.FromFile(config.ConfigFile)

		if err := cliconfig.InitConfig(&pflag.FlagSet{}); err != nil {
			log.Errorf("init config error: ", err)
		}
	})
	return config
}

func parse(str, splitter string) []string {
	var strs []string
	if len(strings.TrimSpace(str)) > 0 {
		ss := strings.Split(str, ",")
		for _, s := range ss {
			strs = append(strs, s)
		}
	}
	return strs
}

func transform(args *Args) *secAction.ArgStruct {
	return &secAction.ArgStruct{
		Func: args.Func, Args: args.Args}
}
