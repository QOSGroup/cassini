package adapter

import (
	"github.com/QOSGroup/qbase/txs"
)

// TxPool cache tx
type TxPool struct {
	size uint32
	pool []*txs.TxQcp
}

// NewTxPool create a tx pool
func NewTxPool(size uint32) *TxPool {
	return &TxPool{
		size: size,
		pool: make([]*txs.TxQcp, size)}
}

// Publish publish a qcp tx
func (p TxPool) Publish(tx *txs.TxQcp) {
	if tx == nil {
		return
	}

}
