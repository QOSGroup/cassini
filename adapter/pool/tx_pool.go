package pool

import (
	"errors"

	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
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
func (p TxPool) Publish(tx *txs.TxQcp) error {
	if tx == nil {
		return errors.New("TxQcp is nil")
	}
	i := uint32(tx.Sequence) % p.size
	p.pool[i] = tx
	log.Tracef("TxPool cached tx: %s", cmn.StringTx(tx))
	return nil
}
