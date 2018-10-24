package adapter

// copy from tendermint/rpc/core/abci.go

import (
	"fmt"
	"strings"

	"github.com/QOSGroup/cassini/log"
	motxs "github.com/QOSGroup/cassini/mock/tx"
	"github.com/QOSGroup/qbase/txs"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// ABCIQuery 交易、交易序号查询。
func ABCIQuery(path string, data cmn.HexBytes, height int64, trusted bool) (*ctypes.ResultABCIQuery, error) {
	if height < 0 {
		log.Errorf("Query sequence error: height [%d] < 0, height must be non-negative", height)
		return nil, fmt.Errorf("height must be non-negative")
	}
	var err error

	cdc := amino.NewCodec()
	ctypes.RegisterAmino(cdc)
	txs.RegisterCodec(cdc)
	cdc.RegisterConcrete(&motxs.TxMock{}, "cassini/mock/txmock", nil)

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

	tx := motxs.NewTxQcpMock("qqs", "qos", height, height)

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
