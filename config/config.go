package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/QOSGroup/cassini/log"
)

// Config 封装配置数据
type Config struct {

	// 消息队列服务配置
	// 如果既没配置Kafka也没配置Nats，则认为配置内部队列模式，仅建议用于测试环境下。

	// Nats 集群配置，以逗号分割
	Nats string `json:"nats,omitempty"`

	// Kafka 集群配置，以逗号分割
	Kafka string `json:"kafka,omitempty"`

	// Mocks 所有需要Mock的服务配置
	Mocks []MockConfig `json:"mocks,omitempty"`

	// 与relay连接的区块链
	QscNames []QscConfig `json:"qscNames,omitempty"`

	EventsListen string `json:"listen,omitempty"`

	EventsQuery string `json:"query,omitempty"`
}

// QscConfig qsc 配置封装
type QscConfig struct {
	Name string `json:"names,omitempty"`
	//链的公钥
	Pubkey string `json:"pubkey,omitempty"`
	//链给relay颁发的证书文件
	Certificate string `json:"certificate,omitempty"`

	//区块链节点地址，多个之间用“，”分割
	NodeAddress string `json:"nodeAddress,omitempty"`
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

// DefaultConfig returns a default configuration for a Tendermint node
func DefaultConfig() *Config {
	return &Config{
		Nats:     "nats://192.168.168.195:4222,nats://192.168.168.195:4223,nats://192.168.168.195:4224",
		Kafka:    "",
		QscNames: DefaultQscConfig(),
	}
}

// DefaultQscConfig 创建默认配置
func DefaultQscConfig() []QscConfig {
	return []QscConfig{
		QscConfig{
			Name: "qos",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "127.0.0.1:26657",
		},
		QscConfig{
			Name: "qqs",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "127.0.0.1:27657",
		},
	}
}

// TestConfig returns a configuration that can be used for testing
func TestConfig() *Config {
	return &Config{
		Nats:     "nats://192.168.168.195:4222",
		Kafka:    "",
		QscNames: TestQscConfig(),
	}
}

// TestQscConfig 创建测试配置
func TestQscConfig() []QscConfig {
	return []QscConfig{
		QscConfig{
			Name: "qos",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "192.168.168.193",
		},
		QscConfig{
			Name: "qqs",
			//链的公钥
			Pubkey: "",
			//链给relay颁发的证书文件
			Certificate: "",
			//区块链节点地址，多个之间用“，”分割
			NodeAddress: "192.168.168.27",
		},
	}
}

// TestQscMockConfig 创建Qsc Mock 测试配置
func TestQscMockConfig() *MockConfig {
	// {
	// 	"mocks":[
	// 	    {
	// 	        "name": "qsc",
	// 	        "rpc": {
	// 		"listen": "tcp://0.0.0.0:27657"
	// 	        }
	// 	    }
	// 	]
	// }
	return &MockConfig{
		Name: "qsc",
		RPC:  RPCConfig{ListenAddress: "tcp://0.0.0.0:27657"},
	}
}
