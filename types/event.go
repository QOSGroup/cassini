package types

import "github.com/tendermint/go-amino"

type CassiniEventDataTx struct {
	From      string `json:"from"` //qsc name 或 qos
	To        string `json:"to"`   //qsc name 或 qos
	Sequence  int64  `json:"sequence"`
	HashBytes []byte `json:"hashBytes"` //TxQcp 做 sha256
}

type Event struct {
	NodeAddress        string `json:"node"` //event 源地址
	CassiniEventDataTx `json:"eventDataTx"` //event 事件
}

func RegisterCassiniTypesAmino(cdc *amino.Codec) {
	//cdc.RegisterInterface((*TMEventData)(nil), nil)
	cdc.RegisterConcrete(CassiniEventDataTx{}, "cassini/event/CassiniEventDataTx", nil)
	cdc.RegisterConcrete(Event{}, "cassini/event/Event", nil)
	cdc.RegisterConcrete(TxQcp{}, "cassini/txqcp/TxQcp", nil)
}
