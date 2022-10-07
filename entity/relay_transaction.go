package entity

type RelayTransaction struct {
	TransferId  string `json:"transferId" gorm:"primaryKey;<-:create"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	Token       string `json:"token"`
	SendAmount  string `json:"sendAmount"`
	RelayAmount string `json:"relayAmount"`
	Nonce       uint64 `json:"nonce"`
	MaxSlippage int32  `json:"maxSlippage"`

	SourceChainId      int32  `json:"srcChainId"`
	DestinationChainId int32  `json:"dstChainId"`
	RelayTransferId    string `json:"relayTransferId"`
	Status             string `json:"status"`
	Result             string `json:"result"`
	RelayToken         string `json:"relayToken"`
	BlockTime          int64  `json:"blockTime"`
}
