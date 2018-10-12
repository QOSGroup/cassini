package mock

// copy from tendermint/rpc/core/routes.go

// Routes 封装接口映射
// TODO: better system than "unsafe" prefix
// NOTE: Amino is registered in rpc/core/types/wire.go.
var Routes = map[string]*RPCFunc{
	// subscribe/unsubscribe are reserved for websocket events.
	"subscribe":       NewWSRPCFunc(Subscribe, "query"),
	"unsubscribe":     NewWSRPCFunc(Unsubscribe, "query"),
	"unsubscribe_all": NewWSRPCFunc(UnsubscribeAll, ""),

	// info API
	"health":               NewRPCFunc(Health, ""),
	"status":               NewRPCFunc(Status, ""),
	"net_info":             NewRPCFunc(NetInfo, ""),
	"blockchain":           NewRPCFunc(BlockchainInfo, "minHeight,maxHeight"),
	"genesis":              NewRPCFunc(Genesis, ""),
	"block":                NewRPCFunc(Block, "height"),
	"block_results":        NewRPCFunc(BlockResults, "height"),
	"commit":               NewRPCFunc(Commit, "height"),
	"tx":                   NewRPCFunc(Tx, "hash,prove"),
	"tx_search":            NewRPCFunc(TxSearch, "query,prove,page,per_page"),
	"validators":           NewRPCFunc(Validators, "height"),
	"dump_consensus_state": NewRPCFunc(DumpConsensusState, ""),
	"consensus_state":      NewRPCFunc(ConsensusState, ""),
	"unconfirmed_txs":      NewRPCFunc(UnconfirmedTxs, "limit"),
	"num_unconfirmed_txs":  NewRPCFunc(NumUnconfirmedTxs, ""),

	// broadcast API
	"broadcast_tx_commit": NewRPCFunc(BroadcastTxCommit, "tx"),
	"broadcast_tx_sync":   NewRPCFunc(BroadcastTxSync, "tx"),
	"broadcast_tx_async":  NewRPCFunc(BroadcastTxAsync, "tx"),

	// abci API
	"abci_query": NewRPCFunc(ABCIQuery, "path,data,height,trusted"),
	"abci_info":  NewRPCFunc(ABCIInfo, ""),
}

// AddUnsafeRoutes 添加非安全接口映射
func AddUnsafeRoutes() {
	// control API
	Routes["dial_seeds"] = NewRPCFunc(UnsafeDialSeeds, "seeds")
	Routes["dial_peers"] = NewRPCFunc(UnsafeDialPeers, "peers,persistent")
	Routes["unsafe_flush_mempool"] = NewRPCFunc(UnsafeFlushMempool, "")

	// profiler API
	Routes["unsafe_start_cpu_profiler"] = NewRPCFunc(UnsafeStartCPUProfiler, "filename")
	Routes["unsafe_stop_cpu_profiler"] = NewRPCFunc(UnsafeStopCPUProfiler, "")
	Routes["unsafe_write_heap_profile"] = NewRPCFunc(UnsafeWriteHeapProfile, "filename")
}
