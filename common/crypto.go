package common

import (
	"encoding/hex"

	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/tendermint/crypto"
)

// Bytes2HexStr 字节数组转换为16 进制字符串
func Bytes2HexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HashTx 交易哈希算法，用以验证交易
func HashTx(tx *txs.TxQcp) []byte {
	return crypto.Sha256(tx.GetSigData())
}
