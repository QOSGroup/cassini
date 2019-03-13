package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	qbtypes "github.com/QOSGroup/qbase/types"
)

// TxMsg implements interface of txs.ITx
type TxMsg struct {
	Data string `json:"data"`
}

var _ txs.ITx = (*TxMsg)(nil)

func newTxMsg(data string) *TxMsg {
	return &TxMsg{Data: data}
}

// ValidateData implements interface of txs.ITx
func (tx *TxMsg) ValidateData(ctx context.Context) error {
	return nil
}

// Exec implements interface of txs.ITx
func (tx *TxMsg) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	return
}

// GetSigner implements interface of txs.ITx
func (tx *TxMsg) GetSigner() []btypes.Address {
	return nil
}

// CalcGas implements interface of txs.ITx
func (tx *TxMsg) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// GetGasPayer implements interface of txs.ITx
func (tx *TxMsg) GetGasPayer() btypes.Address {
	return nil
}

// GetSignData implements interface of txs.ITx
func (tx *TxMsg) GetSignData() []byte {
	return []byte((tx.Data))
}

// NewTxQcp create a new TxQcp
func NewTxQcp(chainID, from, to string, height int64, sequence int64, msg string) *txs.TxQcp {
	itx := newTxMsg(msg)
	txStd := txs.NewTxStd(itx, chainID, qbtypes.NewInt(int64(0)))
	return &txs.TxQcp{
		From:        from,
		To:          to,
		BlockHeight: height,
		TxIndex:     0,
		Sequence:    sequence,
		TxStd:       txStd,
		Extends:     msg}
}
