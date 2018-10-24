package types

// TxQcp 测试用
type TxQcp struct {
	Payload  []byte `json:"payload"` //交易体
	From     string `json:"from"`    //源链ID
	To       string `json:"to"`      //目标链ID
	Sequence int64  `json:"sequence"`
	IsResult bool   `json:"isResult"` //是否为result
	Sign     []byte `json:"sign"`     //对本交易的签名，公链不需要签名
}
