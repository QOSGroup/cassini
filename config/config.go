package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/QOSGroup/cassini/log"
)

// Config 封装配置数据
type Config struct {

	// ConfigFile 配置文件路径，通过命令行设置
	ConfigFile string `json:"config,omitempty"`

	// LogConfigFile 日志配置文件路径，通过命令行设置
	LogConfigFile string `json:"log,omitempty"`

	// Consensus 中继共识配置，为了以后支持更多共识算法，添加次配置
	//    "no" - 不启用共识
	//    默认 - 2/3共识
	Consensus bool `json:"consensus,omitempty"`

	// EventWaitMillitime 交易事件被监听到后需要等待的事件，
	// <=0 不等待
	// >0 等待相应毫秒数
	EventWaitMillitime int64 `json:"eventWaitMillitime,omitempty"`

	// Prikey Cassini relay's private key
	Prikey string `json:"prikey,omitempty"`

	// 消息队列服务配置
	// 如果既没配置Kafka也没配置Nats，则认为配置内部队列模式，仅建议用于测试环境下。

	// Nats 集群配置，以逗号分割
	Nats string `json:"nats,omitempty"`

	// Kafka 集群配置，以逗号分割
	Kafka string `json:"kafka,omitempty"`

	// Mocks 所有需要Mock的服务配置
	Mocks []*MockConfig `json:"mocks,omitempty"`

	// Qscs 与relay连接的区块链相关配置
	Qscs []*QscConfig `json:"qscs,omitempty"`

	// UseEtcd Whether to use etcd or not
	UseEtcd bool `json:"useEtcd,omitempty"`

	// Lock config the lock
	//
	// "etcd://192.168.1.100:2379,192.168.1.101:2379,192.168.1.102:2379"
	Lock string `json:"lock,omitempty"`

	// LockTTL timeout for lock
	//
	// 5 - the lock will be auto-unlock with 5s when lose session
	LockTTL int `json:"lockTTL,omitempty"`

	// EmbedEtcd Whether to start embed etcd or not
	EmbedEtcd bool `json:"embedEtcd,omitempty"`

	// Etcd Embed-etcd config
	Etcd *EtcdConfig `json:"etcd,omitempty"`
}

// QscConfig qsc 配置封装
type QscConfig struct {
	// Name 链名称
	Name string `json:"name,omitempty"`

	// Signature if need sign tx data for this chain
	// true - required
	// false/default - not required
	Signature bool `json:"signature,omitempty"`

	// Pubkey 链的公钥
	Pubkey string `json:"pubkey,omitempty"`

	// Certificate 链给relay颁发的证书文件
	Certificate string `json:"certificate,omitempty"`

	// NodeAddress 区块链节点地址，多个之间用“，”分割
	NodeAddress string `json:"nodes,omitempty"`
}

var conf = &Config{}

// LoadConfig 读取配置数据
func LoadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file: %v error: %v", path, err)
		return nil, err
	}
	return CreateConfig(bytes)
}

// CreateConfig 根据传入数据创建配置
func CreateConfig(bytes []byte) (*Config, error) {
	err := json.Unmarshal(bytes, conf)
	if err != nil {
		log.Errorf("Create config error: %v", err)
		return nil, err
	}
	return conf, nil
}

// GetConfig 获取配置数据
func GetConfig() *Config {
	return conf
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
		Nats:  "nats://127.0.0.1:4222",
		Kafka: "",
		Qscs:  DefaultQscConfig(),
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
		Nats:  "nats://127.0.0.1:4222",
		Kafka: "",
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
