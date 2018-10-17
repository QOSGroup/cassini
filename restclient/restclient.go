package restclient

import (
	"fmt"

	"github.com/QOSGroup/qbase/txs"
	"github.com/pkg/errors"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

// HTTP rpc http 接口调用客户端封装
type HTTP struct {
	remote string
	rpc    *rpcclient.JSONRPCClient
}

// newHTTP 创建rpc http访问客户端 tcp://<host>:<port>
func newHTTP(remote string) *HTTP {
	rc := rpcclient.NewJSONRPCClient(remote)
	cdc := rc.Codec()
	ctypes.RegisterAmino(cdc)
	cdc.RegisterInterface((*txs.ITx)(nil), nil)
	cdc.RegisterConcrete(&txs.QcpTxResult{}, "qbase/txs/QcpTxResult", nil)
	rc.SetCodec(cdc)

	return &HTTP{
		rpc:    rc,
		remote: remote}
}

// ABCIQueryTxQcp 查询交易数据客户端接口实现
func (c *HTTP) ABCIQueryTxQcp(path string, data cmn.HexBytes) (*txs.TxQcp, error) {
	return c.abciQueryTxQcpWithOptions(path, data, client.DefaultABCIQueryOptions)
}

func (c *HTTP) abciQueryTxQcpWithOptions(path string, data cmn.HexBytes, opts client.ABCIQueryOptions) (*txs.TxQcp, error) {
	result := new(txs.TxQcp)
	_, err := c.rpc.Call("abci_query",
		map[string]interface{}{"path": path, "data": data, "height": opts.Height, "trusted": opts.Trusted},
		result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQuery")
	}
	return result, nil
}

// RestClient rpc 远程访问客户端
type RestClient struct {
	// *client.HTTP
	*HTTP
}

// NewRestClient 创建 rpc 远程访问客户端
func NewRestClient(remote string) *RestClient {
	// return &RestClient{HTTP: client.NewHTTP(remote, "")}
	return &RestClient{HTTP: newHTTP(remote)}
}

//GetTxQcp [chainId]/out/sequence //需要输出到"chainId"的qcp tx最大序号
//[chainId]/out/tx_[sequence] //需要输出到"chainId"的每个qcp tx
//[chainId]/in/sequence //已经接受到来自"chainId"的qcp tx最大序号
//[chainId]/in/pubkey //接受来自"chainId"的合法公钥
func (r *RestClient) GetTxQcp(chainID string, sequence int64) (*txs.TxQcp, error) {

	path := chainID + "/out/tx_" + fmt.Sprintf("%d", sequence)

	return r.ABCIQueryTxQcp(path, []byte(""))
}

//GetSequence 查询交易序列号
func (r *RestClient) GetSequence(chainID string, outin string) (int64, error) {

	// path := chainID + "/" + outin + "/sequence"
	// r.ABCIQuery(path, []byte(""))
	var seq int64 = 1
	return seq, nil
}

//PostTxQcp 发布交易
func (r *RestClient) PostTxQcp(chainID string, qcp *txs.TxQcp) error {

	return nil
}
