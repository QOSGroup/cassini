package rpc

// copy from tendermint/rpc/core/routes.go

// Routes rpc request/handler mapper
func (s RequestHandler) Routes() map[string]*RPCFunc {
	// NOTE: Amino codec is registered in types/codec.go.
	var routes = map[string]*RPCFunc{
		// health check API
		"health": NewRPCFunc(s.Health, ""),

		// events API
		"subscribe":       NewWSRPCFunc(s.Subscribe, "query"),
		"unsubscribe":     NewWSRPCFunc(s.Unsubscribe, "query"),
		"unsubscribe_all": NewWSRPCFunc(s.UnsubscribeAll, ""),

		// tx API
		"broadcast_tx_sync": NewRPCFunc(s.BroadcastTxSync, "tx"),
		"abci_query":        NewRPCFunc(s.ABCIQuery, "path,data,height,trusted"),
	}
	return routes
}
