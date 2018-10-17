// Package mock 封装中继外联服务的mock 实现
//
// 实现 qsc 联盟链事件服务接口及交易处理接口
package mock

// copy from tendermint/node/node.go
//       and tendermint/rpc/core/pipe.go

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	types "github.com/tendermint/tendermint/rpc/lib/types"
	stat "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/state/txindex"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// === tendermint/rpc/core/pipe.go

	subscribeTimeout = 5 * time.Second
	defaultPerPage   = 30
	maxPerPage       = 100
)

var (
	// === tendermint/rpc/core/pipe.go

	eventBus  *tmtypes.EventBus // thread safe
	mempool   stat.Mempool
	genDoc    *tmtypes.GenesisDoc // cache the genesis structure
	txIndexer txindex.TxIndexer
)

// StartQscMock 启动单个qsc mock服务
func StartQscMock(mock *config.MockConfig) (context.CancelFunc, error) {
	log.Debug("Start mock: ", mock.Name)

	eventBus = tmtypes.NewEventBus()
	eventBus.Start()
	// defer eventBus.Stop()
	cancel := func() {
		eventBus.Stop()
		log.Debug("Cancel Qsc mock service")
	}

	mempool = MempoolMocker{}

	cdc := amino.NewCodec()
	ctypes.RegisterAmino(cdc)
	cdc.RegisterInterface((*txs.ITx)(nil), nil)
	// cdc.RegisterConcrete(txs.TxStd{}, "qbase/txs/TxStd", nil)
	cdc.RegisterConcrete(txs.QcpTxResult{}, "qbase/txs/QcpTxResult", nil)

	mux := http.NewServeMux()
	wm := NewWebsocketManager(Routes, cdc, EventSubscriber(eventBus))
	mux.HandleFunc("/websocket", wm.WebsocketHandler)
	RegisterRPCFuncs(mux, Routes, cdc)
	// listener
	_, err := StartHTTPServer(
		mock.RPC.ListenAddress,
		mux,
		Config{MaxOpenConnections: 100},
	)
	if err != nil {
		cancel()
		return nil, err
	}

	// 定时发布Tx 事件
	ticker := time.NewTicker(time.Millisecond * 1000)
	var h int64
	var i uint32
	h = 1
	i = 0
	go func() {
		for range ticker.C {
			tx := tmtypes.Tx("abc-just-for-test")
			result := abci.ResponseDeliverTx{
				Data: []byte("mock"),
				Tags: []cmn.KVPair{
					{Key: []byte("qcp.to"), Value: []byte("QOS")},
					{Key: []byte("qcp.from"), Value: []byte("QSC1")},
					{Key: []byte("qcp.sequence"), Value: []byte(fmt.Sprintf("%v", h))},
					{Key: []byte("qcp.hash"), Value: []byte("abc-just-for-test")},
				}}

			err = eventBus.PublishEventTx(tmtypes.EventDataTx{TxResult: tmtypes.TxResult{
				Height: h,
				Index:  i,
				Tx:     tx,
				Result: result,
			}})

			if err != nil {
				log.Error("EventBus publish tx error: ", err)
			}

			log.Debug("EventBus publish tx: ", h, ", ", i)

			h++
			i++
		}
	}()

	// listeners[i] = listener
	return cancel, nil
}

// StartMockRPC 按指定地址端口，以及处理方法，启动 HTTP 服务.
func StartMockRPC(
	listen string,
	handler http.Handler,
	mock *config.MockConfig,
) (listener net.Listener, err error) {
	var proto, addr string
	addrs := strings.SplitN(listen, "://", 2)
	if len(addrs) != 2 {
		return nil, fmt.Errorf(
			"invalid address %s (use tcp://127.0.0.1:27657)", listen)
	}
	proto, addr = addrs[0], addrs[1]

	log.Infof("starting mock rpc: %s", listen)
	listener, err = net.Listen(proto, addr)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to listen: %v error: %v", listen, err)
	}
	// if config.MaxOpenConnections > 0 {
	// 	listener = netutil.LimitListener(listener, config.MaxOpenConnections)
	// }

	go func() {
		err := http.Serve(
			listener,
			recoverHandler(maxBytesHandler{h: handler, n: maxBodyBytes}),
		)
		log.Error("mock rpc stopped! err:", err)
	}()
	return listener, nil
}

// 封装http接口处理方法，记录异常日志。
// 如果封装的方法抛出异常，则恢复记录日志，并返回500状态码
func recoverHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 缓存 ResponseWriter
		rww := &ResponseWriterWrapper{-1, w}
		begin := time.Now()

		// 重设响应头
		origin := r.Header.Get("Origin")
		rww.Header().Set("Access-Control-Allow-Origin", origin)
		rww.Header().Set("Access-Control-Allow-Credentials", "true")
		rww.Header().Set("Access-Control-Expose-Headers", "X-Server-Time")
		rww.Header().Set("X-Server-Time", fmt.Sprintf("%v", begin.Unix()))

		defer func() {
			// Send a 500 error if a panic happens during a handler.
			// Without this, Chrome & Firefox were retrying aborted ajax requests,
			// at least to my localhost.
			if e := recover(); e != nil {

				// If RPCResponse
				if res, ok := e.(types.RPCResponse); ok {
					WriteRPCResponseHTTP(rww, res)
				} else {
					// For the rest,
					log.Error(
						"panic in mock rpc handler", "err", e, "stack",
						string(debug.Stack()),
					)
					rww.WriteHeader(http.StatusInternalServerError)
					WriteRPCResponseHTTP(rww, types.RPCInternalError("", e.(error)))
				}
			}

			// Finally, log.
			durationMS := time.Since(begin).Nanoseconds() / 1000000
			if rww.Status == -1 {
				rww.Status = 200
			}
			log.Info("served mock rpc response",
				"method", r.Method, "url", r.URL,
				"status", rww.Status, "duration", durationMS,
				"remoteAddr", r.RemoteAddr,
			)
		}()

		handler.ServeHTTP(rww, r)
	})
}

// === tendermint/rpc/core/pipe.go

func validatePage(page, perPage, totalCount int) int {
	if perPage < 1 {
		return 1
	}

	pages := ((totalCount - 1) / perPage) + 1
	if page < 1 {
		page = 1
	} else if page > pages {
		page = pages
	}

	return page
}

func validatePerPage(perPage int) int {
	if perPage < 1 || perPage > maxPerPage {
		return defaultPerPage
	}
	return perPage
}
