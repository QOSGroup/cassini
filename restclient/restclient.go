package restclient

import (
	"github.com/QOSGroup/cassini/adapter"
	"github.com/QOSGroup/cassini/log"
	motxs "github.com/QOSGroup/cassini/mock/tx"
	catypes "github.com/QOSGroup/cassini/types"
	bctxs "github.com/QOSGroup/qbase/example/basecoin/tx"
	bctypes "github.com/QOSGroup/qbase/example/basecoin/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/pkg/errors"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// HTTP rpc http 接口调用客户端封装
type HTTP struct {
	remote string
	rpc    *JSONRPCClient
}

// newHTTP 创建rpc http访问客户端 tcp://<host>:<port>
func newHTTP(remote string, cdc *amino.Codec) *HTTP {
	rc := NewJSONRPCClient(remote)
	rc.SetCodec(cdc)

	return &HTTP{
		rpc:    rc,
		remote: remote}
}

// ABCIQuery abci query 标准接口
func (c *HTTP) ABCIQuery(path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return c.abciQueryWithOptions(path, data, client.DefaultABCIQueryOptions)
}

func (c *HTTP) abciQueryWithOptions(path string, data cmn.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.ResultABCIQuery)
	_, err := c.rpc.Call("abci_query",
		map[string]interface{}{"path": path, "data": data, "height": opts.Height, "trusted": opts.Trusted},
		result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQuery")
	}
	return result, nil
}

// BroadcastTxSync 同步交易广播调用接口
func (c *HTTP) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_sync", tx)
}

func (c *HTTP) broadcastTX(route string, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	result := new(ctypes.ResultBroadcastTx)
	_, err := c.rpc.Call(route, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.Wrap(err, route)
	}
	return result, nil
}

// RestClient rpc 远程访问客户端
type RestClient struct {
	// *client.HTTP
	*HTTP
	cdc *amino.Codec
}

// NewRestClient 创建 rpc 远程访问客户端
func NewRestClient(remote string) *RestClient {
	// return &RestClient{HTTP: client.NewHTTP(remote, "")}
	// cdc := app.MakeCodec()
	cdc := amino.NewCodec()
	ctypes.RegisterAmino(cdc)
	txs.RegisterCodec(cdc)
	cdc.RegisterConcrete(&bctypes.AppAccount{}, "basecoin/AppAccount", nil)
	cdc.RegisterConcrete(&bctxs.SendTx{}, "basecoin/SendTx", nil)
	cdc.RegisterConcrete(&motxs.TxMock{}, "cassini/mock/txmock", nil)

	return &RestClient{HTTP: newHTTP(remote, cdc), cdc: cdc}
}

// GetTxQcp 查询指定交易，查询键值参考
func (r *RestClient) GetTxQcp(chainID string, sequence int64) (*txs.TxQcp, error) {
	key := catypes.GetChainOutTxsKey(chainID, sequence)
	result, err := r.ABCIQuery("/store/qcp/key", []byte(key))
	if err != nil || result == nil {
		log.Errorf("Get TxQcp error: %v", err)
		return nil, err
	}

	var tx txs.TxQcp

	if result.Response.GetValue() != nil {
		err = r.cdc.UnmarshalBinaryBare(result.Response.GetValue(), &tx)
		if err != nil {
			log.Errorf("Get TxQcp error: %v", err)
			return nil, err
		}
		log.Debugf("Get TxQcp: %v", adapter.StringTx(&tx))
	}
	return &tx, nil
}

// GetSequence 查询交易序列号
func (r *RestClient) GetSequence(chainID string, outin string) (int64, error) {
	path := "/store/qcp/key"
	key := catypes.GetMaxChainOutSequenceKey(chainID)
	result, err := r.ABCIQuery(path, []byte(key))
	if err != nil {
		log.Errorf("Get sequence error: %v", err)
		return -1, err
	}
	var seq int64
	if result.Response.GetValue() != nil {
		err = r.cdc.UnmarshalBinaryBare(result.Response.GetValue(), &seq)
		if err != nil {
			log.Errorf("Get sequence error when parse: %v", err)
			return -1, err
		}
	}
	log.Debugf("Get sequence: %d", seq)
	return seq, nil
}

//PostTxQcp 广播交易
func (r *RestClient) PostTxQcp(chainID string, qcp *txs.TxQcp) error {
	tx, err := r.cdc.MarshalBinaryBare(qcp)
	if err != nil {
		log.Errorf("Marshal TxQcp error: %v", err)
		return err
	}
	var result *ctypes.ResultBroadcastTx
	result, err = r.BroadcastTxSync(tx)

	if err == nil && result.Code != abci.CodeTypeOK {
		err = errors.New(result.Log)
	}
	if err != nil {
		log.Errorf("Post TxQcp error: %v", err)
		return err
	}
	log.Debugf("Post TxQcp successful - %v", qcp)
	return nil
}
