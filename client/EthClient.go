package client

import (
	"cbridgewrapper/abi"
	"cbridgewrapper/entity"
	"cbridgewrapper/logger"
	"cbridgewrapper/model"
	"cbridgewrapper/utils"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const swapTimeOut = 600000

var clientMap map[int32]*EthClient

func init() {
	logger.Logger.Infof("Initialize client")
	data, err := ioutil.ReadFile("resources/chains.json")
	if err != nil || len(data) == 0 {
		logger.Logger.Fatalf("Invalid chain config")
	}
	clientMap = make(map[int32]*EthClient)

	var chainConfigs []model.ChainConfig
	err = utils.Unmarshal(data, &chainConfigs)
	if err != nil {
		logger.Logger.Fatalf("Parse chain config fail")
	}
	for _, chainConfig := range chainConfigs {
		cl, err := ethclient.Dial(chainConfig.RpcUrl)
		if err != nil {
			logger.Logger.WithField("chainId", chainConfig.ChainId).Fatalf("Can't connect to rpc url %s", chainConfig.RpcUrl)
		}
		Abi, _ := abi.NewToken(common.HexToAddress(chainConfig.Router), cl)
		EthClient := EthClient{
			NetworkId: chainConfig.ChainId,
			WETH:      common.HexToAddress(chainConfig.WEth),
			Abi:       Abi,
			ChainId:   0,
		}
		clientMap[chainConfig.ChainId] = &EthClient
	}
}

func GetEthClient(chainId int32) *EthClient {
	return clientMap[chainId]
}

type EthClient struct {
	NetworkId int32
	WETH      common.Address
	Abi       *abi.Token
	ChainId   int32
}

func (c *EthClient) Swap(transaction entity.RelayTransaction) (*types.Transaction, error) {
	wei := new(big.Int)
	wei.SetString(transaction.RelayAmount, 10)
	return c.Abi.SwapETHForExactTokens(&bind.TransactOpts{}, wei, []common.Address{c.WETH, common.HexToAddress(transaction.RelayToken)}, common.HexToAddress(transaction.Receiver), big.NewInt(utils.MakeTimestamp()+swapTimeOut))
}
