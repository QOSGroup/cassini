// Package mock 封装中继外联服务的mock 实现
//
// 实现 qsc 联盟链事件服务接口及交易处理接口
package mock

// copy from tendermint/node/node.go
//       and tendermint/rpc/core/pipe.go

import (
	"context"
	"time"

	"github.com/QOSGroup/cassini/adapter"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/tendermint/state/txindex"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// === tendermint/rpc/core/pipe.go

	defaultPerPage = 30
	maxPerPage     = 100
)

var (
	// === tendermint/rpc/core/pipe.go

	genDoc    *tmtypes.GenesisDoc // cache the genesis structure
	txIndexer txindex.TxIndexer
)

// StartMock 启动单个 mock服务
func StartMock(mock config.MockConfig) (context.CancelFunc, error) {
	log.Debug("Start mock: ", mock.Name)

	adapter, err := adapter.NewAdapter(mock.Name, "mocktest-id", mock.RPC.ListenAddress, nil, nil)
	cancel := func() {
		adapter.Stop()
	}
	if err != nil {
		return nil, err
	}
	adapter.Start()

	ticker := func(mock *config.MockConfig) {
		log.Debug("ticker: ", mock.Name)
		// 定时发布Tx 事件
		tick := time.NewTicker(time.Millisecond * 1000)
		for range tick.C {
			err = adapter.BroadcastTx(txs.TxQcp{
				From:        mock.Name,
				To:          mock.Name,
				BlockHeight: 1,
				TxIndx:      0,
				Sequence:    0})
			if err != nil {
				log.Error("EventBus publish tx error: ", err)
			}
		}
	}

	go ticker(&mock)

	return cancel, nil
}
