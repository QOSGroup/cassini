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

// BytesInt64 Int64 转换
func BytesInt64(bs []byte) (x int64, err error) {
	buf := bytes.NewBuffer(bs)
	err = binary.Read(buf, binary.BigEndian, &x)
	return x, err
}

// Int64Bytes int64 与 byte 数组转换
func Int64Bytes(in int64) []byte {
	var ret = bytes.NewBuffer([]byte{})
	err := binary.Write(ret, binary.BigEndian, in)
	if err != nil {
		fmt.Printf("Int2Byte error:%s", err.Error())
		return nil
	}

	return ret.Bytes()
}

// GetMaxChainOutSequenceKey 输出队列交易序号查询接口key值组装方法
func GetMaxChainOutSequenceKey(outChain string) string {
	return fmt.Sprintf(outSequenceKey, outChain)
}

// GetChainOutTxsKey 输出队列交易查询接口key值组装方法
func GetChainOutTxsKey(outChain string, sequence int64) string {
	return fmt.Sprintf(outSequenceTxKey, outChain, sequence)
}
