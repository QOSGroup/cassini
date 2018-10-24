package types

import (
	"bytes"
	"encoding/binary"
)

// Bytes2Int64 byte 数组与Int64 转换
func Bytes2Int64(bs []byte) (int64, error) {
	buf := bytes.NewBuffer(bs)
	var x int64
	err := binary.Read(buf, binary.BigEndian, &x)
	return x, err
}
