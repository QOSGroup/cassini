package tx

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	qbtypes "github.com/QOSGroup/qbase/types"
)

// TxMock Mock 交易，实现txs.ITx 接口
type TxMock struct {
	Data string `json:"data"`
}

var _ txs.ITx = (*TxMock)(nil)

func newTxMock(data string) *TxMock {
	return &TxMock{Data: data}
}

// ValidateData Mock 交易，实现txs.ITx 接口
func (tx *TxMock) ValidateData(ctx context.Context) error {
	return nil
}

// Exec Mock 交易，实现txs.ITx 接口
func (tx *TxMock) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	return
}

// GetSigner Mock 交易，实现txs.ITx 接口
func (tx *TxMock) GetSigner() []btypes.Address {
	return nil
}

// CalcGas Mock 交易，实现txs.ITx 接口
func (tx *TxMock) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// GetGasPayer Mock 交易，实现txs.ITx 接口
func (tx *TxMock) GetGasPayer() btypes.Address {
	return nil
}

// GetSignData Mock 交易，实现txs.ITx 接口
func (tx *TxMock) GetSignData() []byte {
	return []byte((tx.Data))
}

// NewTxQcpMock 创建 Mock 交易结构实例
func NewTxQcpMock(from, to string, height int64, sequence int64) *txs.TxQcp {
	itx := newTxMock("Mock Tx by cassini-mock")
	txStd := txs.NewTxStd(itx, "cassini-mock", qbtypes.NewInt(int64(0)))
	return &txs.TxQcp{
		From:        from,
		To:          to,
		BlockHeight: height,
		TxIndex:     0,
		Sequence:    sequence,
		TxStd:       txStd}
}
