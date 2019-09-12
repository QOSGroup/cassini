package common

import (
	"fmt"

	"strconv"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	amino "github.com/tendermint/go-amino"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Transform 将交易转换为交易事件
func Transform(tx *txs.TxQcp) (*tmtypes.EventDataTx, error) {
	result := abcitypes.ResponseDeliverTx{
		Data: []byte("mock")}
	result.Events = append(result.Events, abcitypes.Event{
		Type: types.EventTypeMessage,
		Attributes: []cmn.KVPair{
			{Key: []byte(types.AttributeKeyModule), Value: []byte(qcp.EventModule)},
			{Key: []byte(qcp.To), Value: []byte(tx.To)},
			{Key: []byte(qcp.From), Value: []byte(tx.From)},
			{Key: []byte(qcp.Sequence), Value: []byte(strconv.FormatInt(tx.Sequence, 10))},
			// {Key: []byte(qcp.Hash), Value: []byte(qcp.GenQcpTxHash(tx))},
			{Key: []byte(qcp.Hash), Value: HashTx(tx)},
		}})
	return &tmtypes.EventDataTx{TxResult: tmtypes.TxResult{
		Height: tx.BlockHeight,
		Index:  uint32(tx.TxIndex),
		Tx:     tx.BuildSignatureBytes(),
		Result: result,
	}}, nil
}

// SignTxQcp Sign Tx data for chain
func SignTxQcp(tx *txs.TxQcp, prikey string, cdc *amino.Codec) error {
	//如果密钥是16进制串
	//hex, err := hex.DecodeString(prikey)
	//if err != nil {
	//	return err
	//}
	//var signer ed25519.PrivKeyEd25519
	//cdc.MustUnmarshalBinaryBare(hex, &signer)
	signer, err := UnmarshalKey(prikey)
	if err != nil {
		return err
	}
	tx.Sig.Pubkey = signer.PubKey()
	tx.Sig.Signature, err = tx.SignTx(signer)
	log.Debugf("tx.sig %v", tx.Sig)
	return err
}

// StringTx 将交易转换为字符串，用于日志记录，非完全序列化
func StringTx(tx *txs.TxQcp) string {
	if tx == nil {
		return ""
	}
	return fmt.Sprintf("[%s, %s, %d, %d, %d]", tx.From, tx.To, tx.BlockHeight, tx.TxIndex, tx.Sequence)
}
