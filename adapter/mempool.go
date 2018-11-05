package adapter

import (
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/mempool"
)

// NewMempool create a new mempool
func NewMempool() *mempool.Mempool {
	mc := &config.MempoolConfig{
		WalPath:   "./wal/mempool.wal",
		CacheSize: 1000}
	pool := mempool.NewMempool(mc, nil, 0)
	pool.InitWAL()
	return pool
}
