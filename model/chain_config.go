package model

type ChainConfig struct {
	RpcUrl  string `json:"rpcUrl"`
	WEth    string `json:"weth"`
	Router  string `json:"router"`
	ChainId int32  `json:"chainId"`
	Network string `json:"network"`
}
