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
	txs "github.com/QOSGroup/cassini/mock/tx"
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
	if err != nil {
		return nil, err
	}
	cdc := adapter.GetCodec()
	cdc.RegisterConcrete(&txs.TxMock{}, "cassini/mock/txmock", nil)

	err = adapter.Start()
	if err != nil {
		return nil, err
	}
	cancel := func() {
		adapter.Stop()
	}
	ticker := func(mock *config.MockConfig) {
		log.Debug("ticker: ", mock.Name)
		// 定时发布Tx 事件
		tick := time.NewTicker(time.Millisecond * 1000)
		h := int64(1)
		for range tick.C {
			tx := txs.NewTxQcpMock(mock.Name, mock.To, h, h)
			err = adapter.BroadcastTx(*tx)
			if err != nil {
				log.Error("EventBus publish tx error: ", err)
			}
			h++
		}
	}

	go ticker(&mock)

	return cancel, nil
}
