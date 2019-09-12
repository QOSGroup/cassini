package config

import (
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Config wraps all configure data of cassini
type Config struct {
	// Queue define message queue service type, IP and port addresses.
	// Multiple addresses should be separated by comma.
	// Example:
	//     nats://192.168.1.100:4222,192.168.1.101:4222,192.168.1.102:4222
	// default:
	//	   local
	Queue string `yaml:"queue,omitempty"`

	// Prikey Cassini relay's private key
	Prikey string `yaml:"prikey,omitempty"`

	// Consensus setting the consensus for cassini
	// "no"    - no consensus
	// default - 2/3 consensus
	Consensus bool `yaml:"consensus,omitempty"`

	// EventWaitMillitime 交易事件被监听到后需要等待的事件，
	// <=0 不等待
	// >0 等待相应毫秒数
	EventWaitMillitime int64 `yaml:"eventWaitMillitime,omitempty"`

	// Mocks 所有需要Mock的服务配置
	Mocks []*MockConfig `yaml:"mocks,omitempty"`

	// Qscs 与relay连接的区块链相关配置
	Qscs []*QscConfig `yaml:"qscs,omitempty"`

	// UseEtcd Whether to use etcd or not
	UseEtcd bool `yaml:"useEtcd,omitempty"`

	// Lock config the lock
	//
	// "etcd://192.168.1.100:2379,192.168.1.101:2379,192.168.1.102:2379"
	Lock string `yaml:"lock,omitempty"`

	// LockTTL timeout for lock
	//
	// 5 - the lock will be auto-unlock with 5s when lose session
	LockTTL int64 `yaml:"lockTTL,omitempty"`

	// EmbedEtcd Whether to start embed etcd or not
	EmbedEtcd bool `yaml:"embedEtcd,omitempty"`

	// Etcd Embed-etcd config
	Etcd *EtcdConfig `yaml:"etcd,omitempty"`
}

// QscConfig qsc 配置封装
type QscConfig struct {
	// Name 链名称
	Name string `yaml:"name,omitempty"`

	// Type 链类型
	Type string `yaml:"type,omitempty"`

	// Signature if need sign tx data for this chain
	// true - required
	// false/default - not required
	Signature bool `json:"signature,omitempty"`

	// Pubkey 链的公钥
	Pubkey string `json:"pubkey,omitempty"`

	// Certificate 链给relay颁发的证书文件
	Certificate string `json:"certificate,omitempty"`

	// Nodes 区块链节点地址，多个之间用“，”分割
	Nodes string `yaml:"nodes,omitempty"`
}

var conf = &Config{}

// GetConfig returns the config instance of cassini
func GetConfig() *Config {
	return conf
}

// Load the configure file
func (c *Config) Load() (err error) {
	if err = viper.Unmarshal(c); err != nil {
		return
	}

	var qscs []*QscConfig
	if err = viper.UnmarshalKey("qscs", &qscs); err != nil {
		return
	}
	c.Qscs = qscs

	// var mocks []*MockConfig
	// if err = viper.UnmarshalKey("mocks", &mocks); err != nil {
	// 	return
	// }
	// c.Mocks = mocks

	// TODO ??? whats wrong ???
	// var etcd *EtcdConfig
	// if err = viper.UnmarshalKey("etcd", etcd); err != nil {
	// 	return
	// }
	// c.Etcd = etcd

	var etcd EtcdConfig
	if err = viper.UnmarshalKey("etcd", &etcd); err != nil {
		return
	}
	c.Etcd = &etcd

	return
}

// Parse the configure file
func (c *Config) Parse(bytes []byte) error {
	return yaml.UnmarshalStrict(bytes, c)
}

// GetQscConfig 获取指定 ChainID 的 QSC 配置
func (c *Config) GetQscConfig(chainID string) (qsc QscConfig) {
	if len(c.Qscs) > 0 {
		for _, s := range c.Qscs {
			if strings.EqualFold(chainID, s.Name) {
				qsc = *s
				return
			}
		}
	}
	return
}

// DefaultConfig returns a default configuration for a Tendermint node
func DefaultConfig() *Config {
	return &Config{
		Queue:              "nats://127.0.0.1:4222",
		EventWaitMillitime: 2000,
		Prikey:             "",
		Qscs:               DefaultQscConfig(),
	}
}

// DefaultQscConfig 创建默认配置
func DefaultQscConfig() []*QscConfig {
	return []*QscConfig{
		&QscConfig{
			Name:        "qsc",
			Type:        "qos",
			Signature:   false,             // Whether to sign the transaction
			Pubkey:      "",                // Public key of chain
			Certificate: "",                // Certificate of relayer
			Nodes:       "127.0.0.1:26657", // Chain node address, with "," split between multiple
		},
		&QscConfig{
			Name:        "qos",
			Type:        "qos",
			Signature:   false,
			Pubkey:      "",
			Certificate: "",
			Nodes:       "120.0.0.1:27657,127.0.0.1:28657",
		},
	}
}

// TestConfig returns a configuration that can be used for testing
func TestConfig() *Config {
	return &Config{
		Queue: "nats://127.0.0.1:4222",
		Qscs:  TestQscConfig(),
	}
}

// TestQscConfig 创建测试配置
func TestQscConfig() []*QscConfig {
	return []*QscConfig{
		&QscConfig{
			Name: "qos",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			Nodes: "127.0.0.1",
		},
		&QscConfig{
			Name: "qqs",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			Nodes: "127.0.0.1",
		},
	}
}

// TestQscMockConfig 创建Qsc Mock 测试配置
func TestQscMockConfig() *MockConfig {
	return &MockConfig{
		Name: "qsc",
		RPC:  &RPCConfig{NodeAddress: "0.0.0.0:27657"},
	}
}
