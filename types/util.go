package types

import (
	"bytes"
	"encoding/binary"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/qcp"
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
		log.Infof("Int2Byte error:%s", err.Error())
		return nil
	}

	return ret.Bytes()
}

// Key4OutChainSequence 输出队列交易序号查询接口key值组装方法
func Key4OutChainSequence(outChain string) string {
	return string(qcp.BuildOutSequenceKey(outChain))
}

// Key4InChainSequence 输出队列交易序号查询接口key值组装方法
func Key4InChainSequence(chain string) string {
	return string(qcp.BuildInSequenceKey(chain))
}

// Key4OutChainTx 输出队列交易查询接口key值组装方法
func Key4OutChainTx(outChain string, sequence int64) string {
	return string(qcp.BuildOutSequenceTxKey(outChain, sequence))
}
