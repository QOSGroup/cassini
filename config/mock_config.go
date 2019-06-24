package config

// MockConfig 封装 Mock 配置数据
type MockConfig struct {

	// Name 链名称，mock 时，交易事件的qcp.from
	Name string `yaml:"name,omitempty"`

	// To 目标链名称，mock 时，交易事件的qcp.to
	To string `yaml:"to,omitempty"`

	// Subscribe 交易事件订阅条件，为空时通过To 参数自动拼装
	Subscribe string `yaml:"subscribe,omitempty"`

	// Sequence 交易序列号
	Sequence int64 `yaml:"sequence,omitempty"`

	// RPC RPC相关配置
	RPC *RPCConfig `yaml:"rpc,omitempty"`
}

// RPCConfig 相关配置
type RPCConfig struct {
	// 监听地址端口
	NodeAddress string `yaml:"nodes,omitempty"`
}
