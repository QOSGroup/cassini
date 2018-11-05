// Package adapter 按照 QCP 协议规范封装适配SDK，以实现安全、简便和快捷的接入非 Tendermint 技术栈的区块链。
//
// 为接入链提供交易发布功能，并为中继通信提供标准接口及相关服务（ Http Rpc、Web Socket ）。
package adapter

import (
	"fmt"

	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
	amino "github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Broadcaster 交易广播接口，通过该接口广播的交易即表示需要通过中继跨链提交交易以最终完成交易。
type Broadcaster interface {
	BroadcastTx(tx *txs.TxQcp) error
}

// Receiver 交易接收接口，接收中继从其他接入链发来的跨链交易。
type Receiver interface {
	ReceiveTx(tx *txs.TxQcp) error
}

// HandlerService 中继基础服务封装接口
type HandlerService interface {
	Start() error
	Stop() error
	GetCodec() *amino.Codec
	PublishTx(tx *txs.TxQcp) error
	PublishEvent(e *tmtypes.EventDataTx) error
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
func (b *DefaultBroadcaster) BroadcastTx(tx *txs.TxQcp) (err error) {
	var e *tmtypes.EventDataTx
	e, err = cmn.Transform(tx)
	s := StringTx(tx)
	if err != nil {
		log.Errorf("Transform tx %v error: %v", s, err)
		return
	}
	err = b.adapter.PublishTx(tx)
	if err != nil {
		log.Errorf("Publish tx %v error: %v", s, err)
		return
	}
	err = b.adapter.PublishEvent(e)
	if err != nil {
		log.Errorf("Publish event %v error: %v", s, err)
		return
	}
	log.Debugf("Broadcast tx: sequence[%d] %s", tx.Sequence, s)
	return
}

// StringTx 将交易转换为字符串，用于日志记录，非完全序列化
func StringTx(tx *txs.TxQcp) string {
	if tx == nil {
		return ""
	}
	return fmt.Sprintf("[%v, %v, %v, %v, %v]", tx.From, tx.To, tx.BlockHeight, tx.TxIndex, tx.Sequence)
}
