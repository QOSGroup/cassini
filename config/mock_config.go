package config

// MockConfig 封装 Mock 配置数据
type MockConfig struct {

	// Name 链名称，mock 时，交易事件的qcp.from
	Name string `json:"name,omitempty"`

	// To 目标链名称，mock 时，交易事件的qcp.to
	To string `json:"to,omitempty"`

	// RPC RPC相关配置
	RPC RPCConfig `json:"rpc,omitempty"`
}

// RPCConfig 相关配置
type RPCConfig struct {
	// 监听地址端口
	ListenAddress string `json:"listen,omitempty"`
}
