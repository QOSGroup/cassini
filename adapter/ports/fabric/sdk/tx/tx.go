package tx

// WalletTx fabric wallet Tx
type WalletTx struct {
	Height  string `json:"height,omitempty"`
	Func    string `json:"func,omitempty"`
	Address string `json:"address,omitempty"`
}
