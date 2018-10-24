package common

import "encoding/hex"

func Bytes2HexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
