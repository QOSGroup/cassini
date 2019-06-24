package config

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config wraps all configure data of cassini
type Config struct {

	// ConfigFile is configure file path of cassini
	ConfigFile string

	// LogConfigFile is configure file path of log
	LogConfigFile string `yaml:"log,omitempty"`

	// 消息队列服务配置
	// 如果既没配置Kafka也没配置Nats，则认为配置内部队列模式，仅建议用于测试环境下。

	// Nats 集群配置，以逗号分割
	Nats string `yaml:"nats,omitempty"`

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

	// NodeAddress 区块链节点地址，多个之间用“，”分割
	NodeAddress string `yaml:"nodes,omitempty"`
}

var conf = &Config{}

// GetConfig returns the config instance of cassini
func GetConfig() *Config {
	return conf
}

// Load the configure file
func (c *Config) Load() error {
	bytes, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}
	return c.Parse(bytes)
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
		Nats: "nats://127.0.0.1:4222",
		Qscs: DefaultQscConfig(),
	}
}

// DefaultQscConfig 创建默认配置
func DefaultQscConfig() []*QscConfig {
	return []*QscConfig{
		&QscConfig{
			Name: "qsc",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "127.0.0.1:26657",
		},
		&QscConfig{
			Name: "qos",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "120.0.0.1:27657,127.0.0.1:28657",
		},
	}
}

// TestConfig returns a configuration that can be used for testing
func TestConfig() *Config {
	return &Config{
		Nats: "nats://127.0.0.1:4222",
		Qscs: TestQscConfig(),
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
			NodeAddress: "127.0.0.1",
		},
		&QscConfig{
			Name: "qqs",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "127.0.0.1",
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
