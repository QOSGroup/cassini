package restclient

import (
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/QOSGroup/cassini/types"
)

type RestClient struct {
	http *client.HTTP
}

func  NewRestClient(remote, wsEndpoint string) (r *RestClient) {
	r.http = client.NewHTTP(remote,wsEndpoint)
	return
}


//[chainId]/out/sequence //需要输出到"chainId"的qcp tx最大序号
//[chainId]/out/tx_[sequence] //需要输出到"chainId"的每个qcp tx
//[chainId]/in/sequence //已经接受到来自"chainId"的qcp tx最大序号
//[chainId]/in/pubkey //接受来自"chainId"的合法公钥
func (r *RestClient) GetTxQcp(chainId string,sequence int64) (types.TxQcp, error){

	path := chainId + "/out/tx_"  + string(sequence)
	result, err := r.http.ABCIQuery(path,[]byte(""))
	if err != nil {

	}
	k :=result.Response.Key
	v :=result.Response.Value   //TODO 需要反序列化

	txqcp :=types.TxQcp{To:string(v),From:string(k)} //TODO 组装txqcp结构


	return txqcp,nil
}

//TODO
func (r *RestClient) GetSequence(chainId string,outin string) (int64, error) {

	path := chainId + "/" + outin + "/sequence"
	r.http.ABCIQuery(path,[]byte(""))
	var seq int64=1
	return seq ,nil
}

//TODO
func (r *RestClient) PostTxQcp(chainId string,qcp types.TxQcp)  error {


	return nil
}