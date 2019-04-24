package main

import "C"
import (
	"github.com/QOSGroup/cassini/adapter/ports/fabric/sdk"
)

// ----------------------------------------------------------------------------
// source code for so file generation with "go build " command,
// e.g.
// go build -o fabric.so -buildmode=c-shared main.go
// ----------------------------------------------------------------------------

//export chaincodeInvoke
func chaincodeInvoke(channelID, chaincodeID, args *C.char) *C.char {
	return C.CString(sdk.ChaincodeInvokeByString(C.GoString(channelID),
		C.GoString(chaincodeID), C.GoString(args)))
}

//export chaincodeQuery
func chaincodeQuery(channelID, chaincodeID, args *C.char) *C.char {
	return C.CString(sdk.ChaincodeQueryByString(C.GoString(channelID),
		C.GoString(chaincodeID), C.GoString(args)))
}

//export registerToken
func registerToken(chain, token *C.char) *C.char {
	return C.CString(sdk.RegisterTokenByString(
		C.GoString(chain), C.GoString(token)))
}

func main() {
	// chaincodeInvoke()
}
