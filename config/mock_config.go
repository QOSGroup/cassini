package config

// MockConfig 封装 Mock 配置数据
type MockConfig struct {

	// ChainName 链名称
	Name string `json:"name,omitempty"`

	// RPCConfig RPC相关配置
	RPC RPCConfig `json:"rpc,omitempty"`
}

// RPCConfig 相关配置
type RPCConfig struct {
	// 监听地址端口
	ListenAddress string `json:"listen,omitempty"`
}
