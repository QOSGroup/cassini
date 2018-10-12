package mock

// copy from tendermint/rpc/core/health.go

import (
	"github.com/QOSGroup/cassini/log"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Health Get node health. Returns empty result (200 OK) on success, no response - in
// case of an error.
//
// ```shell
// curl 'localhost:26657/health'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.Health()
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"error": "",
// 	"result": {},
// 	"id": "",
// 	"jsonrpc": "2.0"
// }
// ```
func Health() (*ctypes.ResultHealth, error) {
	log.Debug("RPC call Health")
	return &ctypes.ResultHealth{}, nil
}
