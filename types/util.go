package types

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

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

// ParseHeight parse height in string to int64
func ParseHeight(height string) (int64, error) {
	return strconv.ParseInt(height, 10, 64)
}

// ParseSequence parse sequence in string to int64
func ParseSequence(seq string) (int64, error) {
	return strconv.ParseInt(seq, 10, 64)
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

// ParseAddrs parse protocol and addrs
func ParseAddrs(address string) (protocol string, addrs []string) {
	addrs = strings.SplitN(address, "://", 2)
	if len(addrs) == 2 {
		protocol = addrs[0]
		protocol = strings.TrimSpace(protocol)
		a := addrs[1]
		addrs = strings.Split(a, ",")
	} else {
		protocol, addrs = "", []string{}
	}

	return
}
