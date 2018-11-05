package rpc

// copy from tendermint/rpc/core/abci.go

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/QOSGroup/cassini/log"
	motxs "github.com/QOSGroup/cassini/mock/tx"
	"github.com/QOSGroup/cassini/types"

	"github.com/QOSGroup/qbase/txs"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmttypes "github.com/tendermint/tendermint/types"
)

// ABCIQuery 交易、交易序号查询。
func (s RequestHandler) ABCIQuery(path string, data cmn.HexBytes, height int64, trusted bool) (*ctypes.ResultABCIQuery, error) {
	if height < 0 {
		log.Errorf("Query sequence error: height [%d] < 0, height must be non-negative", height)
		return nil, fmt.Errorf("height must be non-negative")
	}
	var err error

	cdc := types.CreateCompleteCodec()

	key := string(data.Bytes())
	if strings.HasPrefix(key, "sequence/") {
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

	// tx/out/%s/%d
	from, sequence, err := parseTxQueryKey(key)
	if err != nil {
		log.Errorf("Parse tx query key error: ", err)
		return nil, err
	}
	log.Debugf("from: %s, height: %d, sequence: %d", from, height, sequence)
	tx := motxs.NewTxQcpMock(from, "qos", height, sequence)

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

// BroadcastTxSync 广播交易。
func (s RequestHandler) BroadcastTxSync(tx tmttypes.Tx) (*ctypes.ResultBroadcastTx, error) {
	return &ctypes.ResultBroadcastTx{
		Code: abci.CodeTypeOK,
		Hash: tx.Hash(),
	}, nil
}

func parseTxQueryKey(key string) (from string, seq int64, err error) {
	str := strings.Split(key, "/")
	if len(str) < 4 {
		err = fmt.Errorf("Tx query key error: %s", key)
		return
	}
	from = str[2]
	seq, err = strconv.ParseInt(str[3], 10, 64)
	return
}
