// Package adapter 按照 QCP 协议规范封装适配SDK，以实现安全、简便和快捷的接入非 Tendermint 技术栈的区块链。
//
// 为接入链提供交易发布功能，并为中继通信提供标准接口及相关服务（ Http Rpc、Web Socket ）。
package adapter

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
	amino "github.com/tendermint/go-amino"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	stat "github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// === tendermint/rpc/core/pipe.go

	subscribeTimeout = 5 * time.Second
)

// Broadcaster 交易广播接口，通过该接口广播的交易即表示需要通过中继跨链提交交易以最终完成交易。
type Broadcaster interface {
	BroadcastTx(tx txs.TxQcp) error
}

// Receiver 交易接收接口，接收中继从其他接入链发来的跨链交易。
type Receiver interface {
	ReceiveTx(tx txs.TxQcp) error
}

// HandlerService 中继基础服务封装接口
type HandlerService interface {
	Start() error
	Stop() error
	GetCodec() *amino.Codec
	PublishTx(tx txs.TxQcp) error
	PublishEvent(e tmtypes.EventDataTx) error
	CancelTx(tx txs.TxQcp) error
}

// Adapter 适配接口封装，封装交易广播接口和交易接收接口。
//
// 交易广播接口（ Broadcaster ）为调用端检测到交易事件时，由调用端调用;
// 交易接收接口（ Receiver ）为中继适配服务接收到远端中继广播的交易后，由适配服务回调通知调用方接收到远端跨链交易。
type Adapter struct {
	HandlerService
	Broadcaster
	Receiver
}

// NewAdapter 创建新的交易广播器
func NewAdapter(name, id, listenAddr string, r Receiver, b Broadcaster) (*Adapter, error) {
	s, err := NewHandlerService(name, id, listenAddr)
	if err != nil {
		return nil, err
	}
	a := &Adapter{
		HandlerService: s,
		Receiver:       r}
	if b == nil {
		b = &DefaultBroadcaster{adapter: a}
	}
	a.Broadcaster = b
	return a, nil
}

// DefaultBroadcaster 实现内存交易广播器。
//
// 作为接入链跨链交易的缓存以提高查询的相关功能的执行效率。
// 交易在广播器中不会缓存，而是直接转发给中继适配服务。
type DefaultBroadcaster struct {
	adapter *Adapter
}

// BroadcastTx 实现交易广播接口，调用响应的交易及交易事件发布接口，以通过中继跨链提交交易以最终完成交易。
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (b *DefaultBroadcaster) BroadcastTx(tx txs.TxQcp) (err error) {
	var e *tmtypes.EventDataTx
	e, err = Transform(tx)
	s := StringTx(&tx)
	if err != nil {
		log.Errorf("Transform tx %v error: %v", s, err)
		return
	}
	err = b.adapter.PublishTx(tx)
	if err != nil {
		log.Errorf("Publish tx %v error: %v", s, err)
		return
	}
	err = b.adapter.PublishEvent(*e)
	if err != nil {
		log.Errorf("Publish event %v error: %v", s, err)
		ce := b.adapter.CancelTx(tx)
		if ce != nil {
			log.Errorf("Cancel tx %v error: %v", s, err)
		} else {
			log.Debug("Cancel tx: ", s)
		}
		return
	}
	log.Debugf("Broadcast tx: sequence[%d] [%s] ", tx.Sequence, s)
	return
}

// Transform 将交易转换为交易事件
func Transform(tx txs.TxQcp) (*tmtypes.EventDataTx, error) {
	hash := crypto.Sha256(tx.GetSigData())
	result := abcitypes.ResponseDeliverTx{
		Data: []byte("mock"),
		Tags: []cmn.KVPair{
			{Key: []byte("qcp.to"), Value: []byte(tx.To)},
			{Key: []byte("qcp.from"), Value: []byte(tx.From)},
			{Key: []byte("qcp.sequence"), Value: []byte(fmt.Sprintf("%v", tx.Sequence))},
			{Key: []byte("qcp.hash"), Value: hash},
		}}
	return &tmtypes.EventDataTx{TxResult: tmtypes.TxResult{
		Height: tx.BlockHeight,
		Index:  uint32(tx.TxIndex),
		Tx:     tx.GetSigData(),
		Result: result,
	}}, nil
}

// StringTx 将交易转换为字符串，用于日志记录，非完全序列化
func StringTx(tx *txs.TxQcp) string {
	if tx == nil {
		return ""
	}
	return fmt.Sprintf("[%v, %v, %v, %v, %v]", tx.From, tx.To, tx.BlockHeight, tx.TxIndex, tx.Sequence)
}

// DefaultHandlerService 服务管理结构，提供基本的启动，停止等服务管理功能。
//
// 交易事件不会缓存，会直接发送给匹配交易事件订阅条件（ Query ）的远端中继服务；
// 交易会被进行缓存（或持久化存储，取决于选用的交易池服务），交易数、占用内存或缓存时间超过一定量会自动清理缓存。
type DefaultHandlerService struct {
	name          string
	id            string
	listenAddress string
	eventHub      *tmtypes.EventBus
	txPool        stat.Mempool
	mux           *http.ServeMux
	listener      net.Listener
	cdc           *amino.Codec
}

// NewHandlerService 创建新服务管理实例
func NewHandlerService(name, id, listenAddr string) (HandlerService, error) {
	s := &DefaultHandlerService{
		name:          name,
		id:            id,
		listenAddress: listenAddr,
		eventHub:      tmtypes.NewEventBus()}
	err := s.init()

	if err != nil {
		log.Errorf("Init service %v - %v error: %v", s.name, s.id, err)
		return nil, err
	}
	return *s, nil
}

func (s *DefaultHandlerService) init() error {
	s.cdc = amino.NewCodec()
	ctypes.RegisterAmino(s.cdc)
	txs.RegisterCodec(s.cdc)

	s.mux = http.NewServeMux()
	Routes := s.Routes()
	wm := NewWebsocketManager(Routes, s.cdc, EventSubscriber(s.eventHub))
	s.mux.HandleFunc("/websocket", wm.WebsocketHandler)
	RegisterRPCFuncs(s.mux, Routes, s.cdc)
	return nil
}

// Name 获取服务名称，服务名称设置后不允许修改，提供此方法获取服务名称。
func (s DefaultHandlerService) Name() string {
	return s.name
}

// ID 获取服务唯一标识，唯一标识设置后不允许修改，提供此方法获取服务唯一标识。
func (s DefaultHandlerService) ID() string {
	return s.id
}

// Start 启动服务
func (s DefaultHandlerService) Start() (err error) {
	log.Debugf("Start service %v - %v", s.name, s.id)
	s.eventHub.Start()
	s.listener, err = StartHTTPServer(
		s.listenAddress,
		s.mux,
		Config{MaxOpenConnections: 100},
	)
	if err != nil {
		log.Errorf("Start service error: %v", err)
		s.Stop()
		return err
	}
	return nil
}

// Stop 停止服务，释放相关资源
func (s DefaultHandlerService) Stop() error {
	if s.listener != nil {
		s.listener.Close()
	}
	if s.eventHub != nil {
		s.eventHub.Stop()
	}
	log.Debugf("Stop service %v - %v", s.name, s.id)
	return nil
}

// GetCodec 获取amino.Codec 以便　Mock 时修改
func (s DefaultHandlerService) GetCodec() *amino.Codec {
	return s.cdc
}

// PublishTx 发布交易，提供给交易查询
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (s DefaultHandlerService) PublishTx(tx txs.TxQcp) error {
	return nil
}

// PublishEvent 发布交易事件，提供给事件订阅
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (s DefaultHandlerService) PublishEvent(e tmtypes.EventDataTx) error {
	return s.eventHub.PublishEventTx(e)
}

// CancelTx 撤销发布交易
func (s DefaultHandlerService) CancelTx(tx txs.TxQcp) error {
	return nil
}
