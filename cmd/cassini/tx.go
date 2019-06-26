package main

import (
	"context"
	"encoding/hex"
	"fmt"

	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"

	// motxs "github.com/QOSGroup/cassini/mock/tx"
	"github.com/QOSGroup/cassini/restclient"
	// "github.com/QOSGroup/qbase/example/basecoin/app"
	bctxs "github.com/QOSGroup/qbase/example/basecoin/tx"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	bctypes "github.com/QOSGroup/qbase/types"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"

	// ctypes "github.com/tendermint/tendermint/rpc/core/types"
	// cmn "github.com/QOSGroup/cassini/common"
	"github.com/tendermint/tendermint/crypto"
)

// 命令行 tx 命令执行方法
var txHandler = func(conf *config.Config) (context.CancelFunc, error) {
	for _, mockConf := range conf.Mocks {
		// "127.0.0.1:26657"
		log.Debugf("new client: %s, sequence: %d",
			mockConf.RPC.NodeAddress, mockConf.Sequence)
		client := restclient.NewRestClient(mockConf.RPC.NodeAddress)

		// 调用交易查询接口
		tx, err := client.GetTxQcp(mockConf.Name, mockConf.Sequence)
		if err == nil {
			fmt.Printf("Get TxQcp: %s\n", cmn.StringTx(tx))
			fmt.Println("Tx hash: ", cmn.Bytes2HexStr(crypto.Sha256(tx.BuildSignatureBytes())))
			// //////////////
			// hash := cmn.Bytes2HexStr(crypto.Sha256(tx.GetSigData()))
			// fmt.Printf("Tx hash: %s", hash)
		} else {
			fmt.Printf("Get no TxQcp on error: %s\n", err)
		}

		// 调用交易序号查询接口
		var seq int64
		seq, err = client.GetSequence(mockConf.Name, "out")
		if err == nil {
			fmt.Println("Get out sequence: ", seq)
		}
		seq, err = client.GetSequence(mockConf.Name, "in")
		if err == nil {
			fmt.Println("Get in sequence: ", seq)
		}

		// // 调用交易广播接口
		// // -m=qcpTransfer -from=xxx -to=xxx -coin=xxx,xxx -prikey=xxx -nonce=xxx(>=0) -chainid=xxx -qcpprikey=xxx -qcpseq=xxx")
		// sender := "address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay"
		// receiver := "address1srrhd4quypqn0vu5sgrmutgudtnmgm2t2juwya"
		// coin := "qstar"
		// amount := int64(11)
		// prikey := "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc"
		// nonce := int64(1)
		// chainName := "qstar"
		// chainID := "basecoin-chain"
		// qcpPriKey := "0xa3288910405746e29aeec7d5ed56fac138b215e651e3244e6d995f25cc8a74c40dd1ef8d2e8ac876faaa4fb281f17fb9bebb08bc14e016c3a88c6836602ca97595ae32300b"
		// qcpseq := mockConf.Sequence
		// isresult := false

		// // cdc := amino.NewCodec()
		// // ctypes.RegisterAmino(cdc)
		// // txs.RegisterCodec(cdc)
		// // cdc.RegisterConcrete(&bctypes.AppAccount{}, "basecoin/AppAccount", nil)
		// // cdc.RegisterConcrete(&bctxs.SendTx{}, "basecoin/SendTx", nil)
		// cdc := app.MakeCodec()
		// cdc.RegisterConcrete(&motxs.TxMock{}, "cassini/mock/txmock", nil)

		// senderAddr, _ := types.GetAddrFromBech32(sender)
		// receiverAddr, _ := types.GetAddrFromBech32(receiver)

		// txQcp := genQcpTx(cdc, senderAddr, receiverAddr,
		// 	bctypes.BaseCoin{Name: coin, Amount: types.NewInt(amount)},
		// 	prikey, nonce, chainName, chainID, qcpPriKey, qcpseq, isresult)

		// err = client.PostTxQcp("qstar", txQcp)
		// if err == nil {
		// 	fmt.Println(fmt.Sprintf("Post tx is %v", txQcp))
		// }

		// // 调用交易查询接口并计算hash
		// tx, err = client.GetTxQcp("qstar", mockConf.Sequence)
		// if err == nil {
		// 	hash := cmn.Bytes2HexStr(crypto.Sha256(tx.GetSigData()))
		// 	log.Debugf("Tx %s hash: %s", adapter.StringTx(tx), hash)

		// }
	}
	return nil, nil
}

// 生成QCP交易
func genQcpTx(cdc *amino.Codec, sender types.Address, receiver types.Address,
	coin bctypes.BaseCoin, senderPriHex string, nonce int64,
	chainName, chainID string, caPriHex string, qcpseq int64, isresult bool) *txs.TxQcp {
	sendTx := bctxs.NewSendTx(sender, receiver, coin)
	std := txs.NewTxStd(&sendTx, chainID, types.NewInt(int64(0)))
	priHex, _ := hex.DecodeString(senderPriHex[2:])
	var priKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	signData := std.BuildSignatureBytes(nonce, chainID)
	signature, _ := priKey.Sign(signData)
	std.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	tx := txs.NewTxQCP(std, chainName, chainID, qcpseq, 0, 0, isresult, "")
	caHex, _ := hex.DecodeString(caPriHex[2:])
	var caPriKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(caHex, &caPriKey)
	sig, _ := caPriKey.Sign(tx.BuildSignatureBytes())
	tx.Sig.Nonce = qcpseq
	tx.Sig.Signature = sig
	tx.Sig.Pubkey = caPriKey.PubKey()
	return tx
}
