package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	//需要输出到"chainId"的qcp tx最大序号
	outSequenceKey = "sequence/out/%s"
	//需要输出到"chainId"的每个qcp tx
	outSequenceTxKey = "tx/out/%s/%d"
	//已经接受到来自"chainId"的qcp tx最大序号
	inSequenceKey = "sequence/in/%s"
)

// Bytes2Int64 byte 数组与Int64 转换
func Bytes2Int64(bs []byte) (int64, error) {
	buf := bytes.NewBuffer(bs)
	var x int64
	err := binary.Read(buf, binary.BigEndian, &x)
	return x, err
}

// GetMaxChainOutSequenceKey 输出队列交易序号查询接口key值组装方法
func GetMaxChainOutSequenceKey(outChain string) string {
	return fmt.Sprintf(outSequenceKey, outChain)
}

// GetChainOutTxsKey 输出队列交易查询接口key值组装方法
func GetChainOutTxsKey(outChain string, sequence int64) string {
	return fmt.Sprintf(outSequenceTxKey, outChain, sequence)
}
