package rpc

// copy from tendermint/rpc/core/routes.go

// Routes 创建接口路由映射
func (s RequestHandler) Routes() map[string]*RPCFunc {
	// routes 封装接口映射
	// TODO: better system than "unsafe" prefix
	// NOTE: Amino is registered in rpc/core/types/wire.go.
	var routes = map[string]*RPCFunc{
		// subscribe/unsubscribe are reserved for websocket events.
		"subscribe":       NewWSRPCFunc(s.Subscribe, "query"),
		"unsubscribe":     NewWSRPCFunc(s.Unsubscribe, "query"),
		"unsubscribe_all": NewWSRPCFunc(s.UnsubscribeAll, ""),

		// info API
		"health": NewRPCFunc(Health, ""),

		// broadcast API
		"broadcast_tx_sync": NewRPCFunc(s.BroadcastTxSync, "tx"),

		// abci API
		"abci_query": NewRPCFunc(ABCIQuery, "path,data,height,trusted"),
	}
	return routes
}
