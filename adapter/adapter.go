// Package adapter 按照 QCP 协议规范封装适配SDK，以实现安全、简便和快捷的接入非 Tendermint 技术栈的区块链。
//
// 为接入链提供交易发布功能，并为中继通信提供标准接口及相关服务（ Http Rpc、Web Socket ）。
package adapter

import (
	"fmt"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
	stat "github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Broadcaster 交易广播接口，通过该接口广播的交易即表示需要通过中继跨链提交交易以最终完成交易。
type Broadcaster interface {
	BroadcastTx(tx txs.TxQcp) error
}

// Receiver 交易接收接口，接收中继从其他接入链发来的跨链交易。
type Receiver interface {
	ReceiveTx(tx txs.TxQcp) error
}

// Adapter 适配接口封装，封装交易广播接口和交易接收接口。
//
// 交易广播接口（ Broadcaster ）为调用端检测到交易事件时，由调用端调用;
// 交易接收接口（ Receiver ）为中继适配服务接收到远端中继广播的交易后，由适配服务回调通知调用方接收到远端跨链交易。
type Adapter struct {
	Broadcaster
	Receiver
}

// NewAdapter 创建新的交易广播器
func NewAdapter(r Receiver, b Broadcaster) (*Adapter, error) {
	s, err := NewService()
	if err != nil {
		return nil, err
	}
	if b == nil {
		b = &DefaultBroadcaster{srv: s}
	}
	a := &Adapter{
		Broadcaster: b,
		Receiver:    r}
	return a, nil
}

// DefaultBroadcaster 实现内存交易广播器。
//
// 作为接入链跨链交易的缓存以提高查询的相关功能的执行效率。
// 交易在广播器中不会缓存，而是直接转发给中继适配服务。
type DefaultBroadcaster struct {
	srv *Service
}

// BroadcastTx 实现交易广播接口，调用响应的交易及交易事件发布接口，以通过中继跨链提交交易以最终完成交易。
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (b *DefaultBroadcaster) BroadcastTx(tx txs.TxQcp) (err error) {
	var e *tmtypes.EventDataTx
	e, err = Transform(tx)
	s := String(&tx)
	if err != nil {
		log.Errorf("Transform tx %v error: %v", s, err)
		return
	}
	err = b.srv.PublishTx(tx)
	if err != nil {
		log.Errorf("Publish tx %v error: %v", s, err)
		return
	}
	err = b.srv.PublishEvent(*e)
	if err != nil {
		log.Errorf("Publish event %v error: %v", s, err)
		ce := b.srv.CancelTx(tx)
		if ce != nil {
			log.Errorf("Cancel tx %v error: %v", s, err)
		} else {
			log.Debug("Cancel tx: ", s)
		}
		return
	}
	log.Debug("Broadcast tx: ", s)
	return
}

// Transform 将交易转换为交易事件
func Transform(tx txs.TxQcp) (*tmtypes.EventDataTx, error) {
	return nil, nil
}

// String 将交易转换为字符串，用于日志记录，非完全序列化
func String(tx *txs.TxQcp) string {
	if tx == nil {
		return ""
	}
	return fmt.Sprintf("[%v, %v, %v, %v, %v]", tx.From, tx.To, tx.BlockHeight, tx.TxIndx, tx.Sequence)
}

// Service 服务管理结构，提供基本的启动，停止等服务管理功能。
//
// 交易事件不会缓存，会直接发送给匹配交易事件订阅条件（ Query ）的远端中继服务；
// 交易会被进行缓存（或持久化存储，取决于选用的交易池服务），交易数、占用内存或缓存时间超过一定量会自动清理缓存。
type Service struct {
	name string
	id   string
	hub  *tmtypes.EventBus
	pool *stat.Mempool
}

// NewService 创建新服务管理实例
func NewService() (*Service, error) {
	s := &Service{
		hub: tmtypes.NewEventBus()}
	return s, nil
}

// Name 获取服务名称，服务名称设置后不允许修改，提供此方法获取服务名称。
func (s Service) Name() string {
	return s.name
}

// ID 获取服务唯一标识，唯一标识设置后不允许修改，提供此方法获取服务唯一标识。
func (s Service) ID() string {
	return s.id
}

// Start 启动服务
func (s Service) Start() error {
	return nil
}

// Stop 停止服务，释放相关资源
func (s Service) Stop() error {
	return nil
}

// PublishTx 发布交易，提供给交易查询
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (s Service) PublishTx(tx txs.TxQcp) error {
	return nil
}

// PublishEvent 发布交易事件，提供给事件订阅
//
// 因为按照 QCP 协议规范定义，中继都是在接收到交易事件后查询交易数据，因此应保证先调用发布交易接口，然后再调用发布事件接口。
func (s Service) PublishEvent(e tmtypes.EventDataTx) error {
	return nil
}

// CancelTx 撤销发布交易
func (s Service) CancelTx(tx txs.TxQcp) error {
	return nil
}
