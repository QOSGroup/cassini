package adapter

// copy from tendermint/rpc/core/routes.go

// Routes 创建接口路由映射
func (s DefaultHandlerService) Routes() map[string]*RPCFunc {
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
		"broadcast_tx_commit": NewRPCFunc(s.BroadcastTxCommit, "tx"),
		"broadcast_tx_sync":   NewRPCFunc(s.BroadcastTxSync, "tx"),
		"broadcast_tx_async":  NewRPCFunc(s.BroadcastTxAsync, "tx"),

		// abci API
		"abci_query": NewRPCFunc(ABCIQuery, "path,data,height,trusted"),
		"abci_info":  NewRPCFunc(ABCIInfo, ""),
	}
	return routes
}
