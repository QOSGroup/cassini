package adapter

// copy from tendermint/rpc/core/abci.go

import (
	"errors"
	"fmt"
	"strings"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// ABCIQuery Query the application for some information.
//
// ```shell
// curl 'localhost:26657/abci_query?path=""&data="abcd"&trusted=false'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.ABCIQuery("", "abcd", true)
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"error": "",
// 	"result": {
// 		"response": {
// 			"log": "exists",
// 			"height": 0,
// 			"proof": "010114FED0DAD959F36091AD761C922ABA3CBF1D8349990101020103011406AA2262E2F448242DF2C2607C3CDC705313EE3B0001149D16177BC71E445476174622EA559715C293740C",
// 			"value": "61626364",
// 			"key": "61626364",
// 			"index": -1,
// 			"code": 0
// 		}
// 	},
// 	"id": "",
// 	"jsonrpc": "2.0"
// }
// ```
//
// ### Query Parameters
//
// | Parameter | Type   | Default | Required | Description                                    |
// |-----------+--------+---------+----------+------------------------------------------------|
// | path      | string | false   | false    | Path to the data ("/a/b/c")                    |
// | data      | []byte | false   | true     | Data                                           |
// | height    | int64 | 0       | false    | Height (0 means latest)                        |
// | trusted   | bool   | false   | false    | Does not include a proof of the data inclusion |
func ABCIQuery(path string, data cmn.HexBytes, height int64, trusted bool) (*ctypes.ResultABCIQuery, error) {
	if height < 0 {
		log.Errorf("Query sequence error: height [%d] < 0, height must be non-negative", height)
		return nil, fmt.Errorf("height must be non-negative")
	}
	var err error

	// tr := txs.NewQcpTxResult(int64(abci.CodeTypeOK), &[]cmn.KVPair{}, 0, types.NewInt(1111111111), "ok")

	// tstd := txs.NewTxStd(tr, "QOS", types.NewInt(999999999))

	// tx := &txs.TxQcp{
	// 	From:        "QOS",
	// 	To:          "QSC",
	// 	BlockHeight: height,
	// 	TxIndx:      -1,
	// 	Sequence:    0,
	// 	Payload:     *tstd}
	// // 仅作为调试接口使用，实际通信数据结构还不确定。

	// return tx, nil

	cdc := amino.NewCodec()
	ctypes.RegisterAmino(cdc)
	// cdc.RegisterInterface((*txs.ITx)(nil), nil)
	// cdc.RegisterConcrete(&txs.QcpTxResult{}, "qbase/txs/QcpTxResult", nil)
	txs.RegisterCodec(cdc)

	// key := "[qstar]/out/sequence"
	key := string(data.Bytes())
	if strings.HasSuffix(key, "/sequence") {
		seq := int32(3)

		var bytes []byte
		bytes, err = cdc.MarshalBinaryBare(seq)
		if err != nil {
			log.Errorf("Query sequence error: ", err)
			return nil, err
		}

		if err != nil {
			log.Errorf("Query sequence error: ", err)
			return nil, err
		}

		resQuery := &abci.ResponseQuery{
			Log:    "ok: query sequence",
			Height: height,
			Key:    []byte(key),
			Value:  bytes}

		log.Info("ABCIQuery", "path", path, "data", data, "height", height, "result", resQuery)
		return &ctypes.ResultABCIQuery{Response: *resQuery}, nil
	}

	tstd := &txs.TxStd{
		ITx:       nil,
		Signature: nil,
		ChainID:   "QOS",
		MaxGas:    types.NewInt(999999999)}

	tx := &txs.TxQcp{
		From:        "qstar",
		To:          "qos",
		BlockHeight: height,
		TxIndx:      -1,
		Sequence:    0,
		Payload:     *tstd}

	var bytes []byte
	bytes, err = cdc.MarshalBinaryBare(tx)
	if err != nil {
		log.Errorf("Query TxQcp error: ", err)
		return nil, err
	}

	if err != nil {
		log.Errorf("Query TxQcp error: ", err)
		return nil, err
	}

	resQuery := &abci.ResponseQuery{
		Log:    "ok: query TxQcp",
		Height: height,
		Key:    []byte(key),
		Value:  bytes}

	var seq txs.TxQcp
	cdc.UnmarshalBinaryBare(bytes, &seq)
	log.Debugf("Unmarshal seq: %s", seq.From)

	return &ctypes.ResultABCIQuery{Response: *resQuery}, nil
}

// ABCIInfo Get some info about the application.
//
// ```shell
// curl 'localhost:26657/abci_info'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// info, err := client.ABCIInfo()
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"error": "",
// 	"result": {
// 		"response": {
// 			"data": "{\"size\":3}"
// 		}
// 	},
// 	"id": "",
// 	"jsonrpc": "2.0"
// }
// ```
func ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	// resInfo, err := proxyAppQuery.InfoSync(abci.RequestInfo{Version: version.Version})
	// if err != nil {
	// 	return nil, err
	// }
	// return &ctypes.ResultABCIInfo{*resInfo}, nil
	return nil, errors.New("not implemented yet")
}
