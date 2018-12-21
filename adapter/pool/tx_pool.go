package pool

import (
	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
)

// Compensator TxPool would try to get overwritten tx by calling the compensator.
type Compensator func(sequence int64) (*txs.TxQcp, error)

// TxPool cache tx
// No need to clean tx-cache
type TxPool struct {
	size        uint32
	pool        []*txs.TxQcp
	compensator Compensator
}

// NewTxPool create a tx pool
func NewTxPool(size uint32, compensator Compensator) *TxPool {
	return &TxPool{
		size:        size,
		pool:        make([]*txs.TxQcp, size),
		compensator: compensator}
}

// Put cache a qcp tx
func (p TxPool) Put(tx *txs.TxQcp) (err error) {
	if tx == nil {
		return
	}
	i := uint32(tx.Sequence) % p.size
	p.pool[i] = tx
	log.Tracef("TxPool put tx: %d: %s", i, cmn.StringTx(tx))
	return
}

// Get get a qcp tx
func (p TxPool) Get(sequence int64) (tx *txs.TxQcp, err error) {
	i := uint32(sequence) % p.size
	tx = p.pool[i]
	if tx != nil {
		log.Tracef("TxPool get tx: %d: %s", sequence, cmn.StringTx(tx))
		if tx.Sequence == sequence {
			return
		}
		tx = nil
	}
	if p.compensator != nil {
		tx, err = p.compensator(sequence)
		if tx != nil {
			log.Tracef("TxPool compensate tx: %d: %s", sequence, cmn.StringTx(tx))
		}
	}
	return
}
