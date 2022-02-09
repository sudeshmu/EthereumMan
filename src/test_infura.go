package main

import (
	"context"
	"fmt"
	"strings"

	"math/big"

	"etherman/src/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type LogTransferData struct {
	SenderAddress   string
	RecieverAddress string
	ContractAddress string
	SenderBalance   string
	RecieverBalance string
	Tokens          string
}

func main() {

	contractAbi, err := abi.JSON(strings.NewReader(string(utils.ABI())))
	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(contractAbi)

	endpoint := "wss://mainnet.infura.io/ws/v3/dce4801c430e46e6a041c7c9fa01edc8"
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		// // log.Fatal(err)
		zap.S().Fatalw("client dial error ", "message", err.Error())
		// zap.S().Fatalw("client dial error ", "message", err.Error())
	}
	hash := common.HexToAddress("0x809060c546574A0d7084eB17E14373E0628feb60")
	resultChan := make(chan types.Log)
	query := ethereum.FilterQuery{Addresses: []common.Address{hash}}
	subscription, err := client.SubscribeFilterLogs(context.Background(), query, resultChan) // creating go eth clinet
	fmt.Print("test1")
	if err != nil {
		// log.Fatal(err)
		zap.S().Fatalw("subscription error", "message", err.Error())
	}

	logTransferSig2 := []byte("TransferSingle(address,address,address,uint256,uint256)")
	logTransferSigHash2 := crypto.Keccak256Hash(logTransferSig2)

	logBatchTransferSig2 := []byte("TransferBatch(address,address,address,uint256[],uint256[])")
	logBatchTransferSigHash2 := crypto.Keccak256Hash(logBatchTransferSig2)

	for {
		select {
		case sErr := <-subscription.Err():
			fmt.Println("Client - Saw error: ", sErr)
		case data := <-resultChan:
			fmt.Printf("New event")
			fmt.Printf("%v", data)
			fmt.Printf("Topics type %s", data.Topics[0].Hex())
			zap.S().Infow("New Transfer Event ", "address", data.BlockNumber)

			if data.Topics[0].Hex() == logTransferSigHash2.Hex() {
				println("Transfer 2")
			}
			if data.Topics[0].Hex() == logBatchTransferSigHash2.Hex() {
				println("Batch T 2")
			}
			var logTransferData = LogTransferData{}

			fmt.Print(logTransferData)
		}

	}
}
