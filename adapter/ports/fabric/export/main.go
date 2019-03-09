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
func chaincodeInvoke(chainID, chaincodeID, args *C.char) *C.char {
	return C.CString(sdk.ChaincodeInvoke(C.GoString(chainID),
		C.GoString(chaincodeID), C.GoString(args)))
}

func main() {

}
