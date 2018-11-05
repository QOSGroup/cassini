package common

import (
	"encoding/hex"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	amino "github.com/tendermint/go-amino"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Transform 将交易转换为交易事件
func Transform(tx *txs.TxQcp) (*tmtypes.EventDataTx, error) {
	hash := HashTx(tx)
	result := abcitypes.ResponseDeliverTx{
		Data: []byte("mock"),
		Tags: []cmn.KVPair{
			{Key: []byte(qcp.QcpTo), Value: []byte(tx.To)},
			{Key: []byte(qcp.QcpFrom), Value: []byte(tx.From)},
			{Key: []byte(qcp.QcpSequence), Value: types.Int64Bytes(tx.Sequence)},
			{Key: []byte(qcp.QcpHash), Value: hash},
		}}
	return &tmtypes.EventDataTx{TxResult: tmtypes.TxResult{
		Height: tx.BlockHeight,
		Index:  uint32(tx.TxIndex),
		Tx:     tx.GetSigData(),
		Result: result,
	}}, nil
}

// SignTxQcp Sign Tx data for chain
func SignTxQcp(tx *txs.TxQcp, prikey string, cdc *amino.Codec) error {
	hex, err := hex.DecodeString(prikey)
	if err != nil {
		return err
	}
	var signer ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(hex, &signer)

	tx.Sig.Pubkey = signer.PubKey()
	tx.Sig.Signature, err = tx.SignTx(signer)
	log.Infof("tx.sig %v", tx.Sig)
	return err
}
