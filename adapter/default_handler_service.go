package adapter

import (
	"net"
	"net/http"

	"github.com/QOSGroup/cassini/adapter/rpc"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/mempool"
	tmtypes "github.com/tendermint/tendermint/types"
)

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

// DefaultHandlerService 服务管理结构，提供基本的启动，停止等服务管理功能。
//
// 交易事件不会缓存，会直接发送给匹配交易事件订阅条件（ Query ）的远端中继服务；
// 交易会被进行缓存（或持久化存储，取决于选用的交易池服务），交易数、占用内存或缓存时间超过一定量会自动清理缓存。
type DefaultHandlerService struct {
	name          string
	id            string
	listenAddress string
	eventHub      *tmtypes.EventBus
	txPool        *mempool.Mempool
	mux           *http.ServeMux
	listener      net.Listener
	cdc           *amino.Codec
	handler       *rpc.RequestHandler
}

// Init init default handler service
func (s *DefaultHandlerService) init() error {
	s.cdc = types.CreateCompleteCodec()

	s.mux = http.NewServeMux()
	s.handler = &rpc.RequestHandler{
		EventHub: s.eventHub}
	Routes := s.handler.Routes()
	wm := rpc.NewWebsocketManager(Routes, s.cdc, rpc.EventSubscriber(s.eventHub))
	s.mux.HandleFunc("/websocket", wm.WebsocketHandler)
	rpc.RegisterRPCFuncs(s.mux, Routes, s.cdc)
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
	s.listener, err = rpc.StartHTTPServer(
		s.listenAddress,
		s.mux,
		rpc.Config{MaxOpenConnections: 100},
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
func (s DefaultHandlerService) PublishTx(tx *txs.TxQcp) error {

	return nil
}

// PublishEvent 发布交易事件，提供给事件订阅
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (s DefaultHandlerService) PublishEvent(e *tmtypes.EventDataTx) error {
	return s.eventHub.PublishEventTx(*e)
}
