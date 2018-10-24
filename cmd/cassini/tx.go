package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
	motxs "github.com/QOSGroup/cassini/mock/tx"
	"github.com/QOSGroup/qbase/example/basecoin/app"
	bctxs "github.com/QOSGroup/qbase/example/basecoin/tx"
	bctypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	// ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// 命令行 tx 命令执行方法
var txHandler = func(conf *config.Config) (context.CancelFunc, error) {
	log.Debug("tx")
	client := restclient.NewRestClient("tcp://192.168.168.196:26657")

	// 调用交易查询接口
	tx, err := client.GetTxQcp("qstar", 1)
	if err == nil {
		fmt.Printf("Get TxQcp chain: %s to: %s\n", tx.From, tx.To)
	}

	// 调用交易序号查询接口
	var seq int64
	seq, err = client.GetSequence("qstar", "out")
	if err == nil {
		fmt.Println("Get sequence: ", seq)
	}

	// 调用交易广播接口
	// -m=qcpTransfer -from=xxx -to=xxx -coin=xxx,xxx -prikey=xxx -nonce=xxx(>=0) -chainid=xxx -qcpprikey=xxx -qcpseq=xxx")
	sender := "address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay"
	receiver := "address1srrhd4quypqn0vu5sgrmutgudtnmgm2t2juwya"
	coin := "qstar"
	amount := int64(11)
	prikey := "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc"
	nonce := int64(1)
	chainID := "qstar"
	qcpPriKey := "0xa3288910405746e29aeec7d5ed56fac138b215e651e3244e6d995f25cc8a74c40dd1ef8d2e8ac876faaa4fb281f17fb9bebb08bc14e016c3a88c6836602ca97595ae32300b"
	qcpseq := int64(8)
	isresult := false

	// cdc := amino.NewCodec()
	// ctypes.RegisterAmino(cdc)
	// txs.RegisterCodec(cdc)
	// cdc.RegisterConcrete(&bctypes.AppAccount{}, "basecoin/AppAccount", nil)
	// cdc.RegisterConcrete(&bctxs.SendTx{}, "basecoin/SendTx", nil)
	cdc := app.MakeCodec()
	cdc.RegisterConcrete(&motxs.TxMock{}, "cassini/mock/txmock", nil)

	senderAddr, _ := types.GetAddrFromBech32(sender)
	receiverAddr, _ := types.GetAddrFromBech32(receiver)

	txQcp := genQcpTx(cdc, senderAddr, receiverAddr,
		bctypes.BaseCoin{Name: coin, Amount: types.NewInt(amount)},
		prikey, nonce, chainID, qcpPriKey, qcpseq, isresult)

	fmt.Printf("TxQcp: %v\n", txQcp)

	err = client.PostTxQcp("qstar", txQcp)
	if err == nil {
		fmt.Println(fmt.Sprintf("Post tx is %v", txQcp))
	}
	return nil, nil
}

// 生成QCP交易
func genQcpTx(cdc *amino.Codec, sender types.Address, receiver types.Address,
	coin bctypes.BaseCoin, senderPriHex string, nonce int64,
	chainID string, caPriHex string, qcpseq int64, isresult bool) *txs.TxQcp {
	sendTx := bctxs.NewSendTx(sender, receiver, coin)
	std := txs.NewTxStd(&sendTx, "basecoin-chain", types.NewInt(int64(0)))
	priHex, _ := hex.DecodeString(senderPriHex[2:])
	var priKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	signData := append(std.GetSignData(), types.Int2Byte(int64(nonce))...)
	signature, _ := priKey.Sign(signData)
	std.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	tx := txs.NewTxQCP(std, chainID, "basecoin-chain", qcpseq, 0, 0, isresult)
	caHex, _ := hex.DecodeString(caPriHex[2:])
	var caPriKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(caHex, &caPriKey)
	sig, _ := caPriKey.Sign(tx.GetSigData())
	tx.Sig.Nonce = qcpseq
	tx.Sig.Signature = sig
	tx.Sig.Pubkey = caPriKey.PubKey()
	return tx
}
