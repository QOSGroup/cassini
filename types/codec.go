package types

import (
	motxs "github.com/QOSGroup/cassini/mock/tx"
	bcapp "github.com/QOSGroup/qbase/example/basecoin/app"
	qosapp "github.com/QOSGroup/qos/app"
	"github.com/tendermint/go-amino"
	tmttypes "github.com/tendermint/tendermint/types"
)

// CreateCompleteCodec 创建完整（包括：联盟链，公链，中继）amino编码器
func CreateCompleteCodec() *amino.Codec {

	// qos cdc
	cdc := qosapp.MakeCodec()

	// tedermint cdc
	// ctypes.RegisterAmino(cdc)
	// ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmttypes.RegisterEventDatas(cdc)
	tmttypes.RegisterEvidences(cdc)

	// qbase cdc
	bcapp.RegisterCodec(cdc)

	// cassini cdc
	RegisterCassiniTypesAmino(cdc)
	return cdc
}

// RegisterCassiniTypesAmino 注册中继自定义类型
func RegisterCassiniTypesAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(CassiniEventDataTx{}, "cassini/event/CassiniEventDataTx", nil)
	cdc.RegisterConcrete(Event{}, "cassini/event/Event", nil)
	cdc.RegisterConcrete(&motxs.TxMock{}, "cassini/mock/TxMock", nil)
}
