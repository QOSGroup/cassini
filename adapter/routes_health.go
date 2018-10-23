package adapter

// copy from tendermint/rpc/core/health.go

import (
	"github.com/QOSGroup/cassini/log"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Health 健康检查。
func Health() (*ctypes.ResultHealth, error) {
	log.Debug("RPC call Health")
	return &ctypes.ResultHealth{}, nil
}
