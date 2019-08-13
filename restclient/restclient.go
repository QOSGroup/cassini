package restclient

import (
	"strings"

	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
	catypes "github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/pkg/errors"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	tcmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// HTTP rpc http 接口调用客户端封装
type HTTP struct {
	remote string
	rpc    *JSONRPCClient
}

const ERR_emptyqcp = "empty qcp transaction"

// newHTTP 创建rpc http访问客户端 tcp://<host>:<port>
func newHTTP(remote string, cdc *amino.Codec) *HTTP {
	rc := NewJSONRPCClient(remote)
	rc.SetCodec(cdc)

	return &HTTP{
		rpc:    rc,
		remote: remote}
}

// ABCIQuery abci query 标准接口
func (c *HTTP) ABCIQuery(path string, data tcmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return c.abciQueryWithOptions(path, data, client.DefaultABCIQueryOptions)
}

func (c *HTTP) abciQueryWithOptions(path string, data tcmn.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.ResultABCIQuery)
	_, err := c.rpc.Call("abci_query",
		map[string]interface{}{"path": path, "data": data, "height": opts.Height, "prove": opts.Prove},
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
	cdc := catypes.CreateCompleteCodec()
	return &RestClient{HTTP: newHTTP(remote, cdc), cdc: cdc}
}

// GetTxQcp 查询指定交易，查询键值参考
func (r *RestClient) GetTxQcp(chainID string, sequence int64) (*txs.TxQcp, error) {
	key := catypes.Key4OutChainTx(chainID, sequence)
	result, err := r.ABCIQuery("/store/qcp/key", []byte(key))
	if err != nil || result == nil {
		log.Errorf("remote [%s] Get TxQcp error: %v", r.remote, err) //TODO 节点0 down 在此报错
		return nil, err
	}

	var tx txs.TxQcp

	if result.Response.GetValue() != nil {
		err = r.cdc.UnmarshalBinaryBare(result.Response.GetValue(), &tx)
		if err != nil {
			log.Errorf("remote [%s] Get TxQcp error: %v", r.remote, err)
			return nil, err
		}

		log.Debugf("remote [%s] Get TxQcp: %v", r.remote, cmn.StringTx(&tx))
		return &tx, nil
	}
	//log.Errorf("empty qcp transaction chainid [%s] sequence [#%d]", chainID, sequence)
	return nil, errors.New(ERR_emptyqcp)
}

// GetSequence 查询交易序列号
func (r *RestClient) GetSequence(chainID string, outin string) (int64, error) {
	path := "/store/qcp/key"
	var key string
	if strings.EqualFold("in", outin) {
		key = catypes.Key4InChainSequence(chainID)
	} else {
		key = catypes.Key4OutChainSequence(chainID)
	}
	result, err := r.ABCIQuery(path, []byte(key))
	if err != nil {
		//log.Warnf("remote [%s] Get sequence error: %v", r.remote, err)
		return -1, err
	}
	var seq int64
	if result.Response.GetValue() != nil {
		err = r.cdc.UnmarshalBinaryBare(result.Response.GetValue(), &seq)
		if err != nil {
			//log.Warnf("remote [%s] Get sequence error when parse: %v", r.remote, err)
			return -1, err
		}
	}
	//log.Debugf("remote [%s] get sequence: [#%d]", r.remote, seq)
	return seq, nil
}

//PostTxQcp 广播交易
func (r *RestClient) PostTxQcp(chainID string, qcp *txs.TxQcp) error {
	tx, err := r.cdc.MarshalBinaryBare(qcp)
	if err != nil {
		log.Errorf("remote %s[%s] Marshal TxQcp error: %v", chainID, r.remote, err)
		return err
	}
	var result *ctypes.ResultBroadcastTx
	result, err = r.BroadcastTxSync(tx)

	if err == nil && result.Code != abci.CodeTypeOK {
		log.Warnf("abci result code: %d, log: %s", result.Code, result.Log)
		json, err := r.cdc.MarshalJSON(qcp)
		if err != nil {
			log.Errorf("TxQcp mashal json error: %v", err)
		}
		log.Warnf("TxQcp: %s", string(json))
		err = errors.New(result.Log)
	}
	if err != nil {
		log.Errorf("remote %s[%s] Post TxQcp error: %v", chainID, r.remote, err)
		log.Infof("qcp: %v", qcp)
		return err
	}
	log.Debugf("remote %s[%s] Post TxQcp successful. %v", chainID, r.remote, qcp)
	return nil
}
