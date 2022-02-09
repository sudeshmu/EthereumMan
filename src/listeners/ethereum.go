package listeners

import (
	"context"
	"math/big"
	"strings"
	"time"

	"etherman/src/contracts"
	log "etherman/src/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LogTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Tokens   *big.Int
}

type LogTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Tokens   []*big.Int
}

func Listen() {

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/56ab5afaf4d9451da8a2a72225d02aba")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(configs.SmartContractAdd())
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	instance, err := contracts.NewToken(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("TransferSingle(address,address,address,uint256,uint256)")
	logTransferBatchSig := []byte("TransferBatch(address,address,address,uint256[],uint256[])")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logTransferBatchSigHash := crypto.Keccak256Hash(logTransferBatchSig)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():

				var transferSingleEvent LogTransferSingle

				err := contractAbi.UnpackIntoInterface(&transferSingleEvent, "TransferSingle", vLog.Data)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatal("could not unpack contract Abi")
				}

				transferSingleEvent.From = common.HexToAddress(vLog.Topics[2].Hex())
				transferSingleEvent.To = common.HexToAddress(vLog.Topics[3].Hex())

				// find the balance of sender
				senderBal, err := instance.BalanceOf(&bind.CallOpts{}, transferSingleEvent.From, transferSingleEvent.Id)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Info("could not obtain sender balance")
				}

				// find the balance of the receiver
				receiverBal, err := instance.BalanceOf(&bind.CallOpts{}, transferSingleEvent.To, transferSingleEvent.Id)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Info("could not obtain receivers balance")
				}

				newTransaction := models.Transaction{
					TokenId:     transferSingleEvent.Id.String(),
					Sender:      transferSingleEvent.From.Hex(),
					Receiver:    transferSingleEvent.To.Hex(),
					SenderBal:   senderBal.String(),
					ReceiverBal: receiverBal.String(),
					Token:       transferSingleEvent.Tokens.String(),
					CreatedAt:   time.Now(),
				}

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

				defer cancel()
				result, err := transactionCollection.InsertOne(ctx, newTransaction)

				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Info("error saving to mongodb")
				}

				log.WithFields(log.Fields{
					"From":        transferSingleEvent.From.Hex(),
					"To":          transferSingleEvent.To.Hex(),
					"Tokens":      transferSingleEvent.Tokens,
					"SenderBal":   senderBal.String(),
					"ReceiverBal": receiverBal.String(),
					"result":      result,
				}).Info("new transaction parsed")

			// Didnt completely handle the batch transfer logging
			case logTransferBatchSigHash.Hex():

				var transferBatchEvent LogTransferBatch

				err := contractAbi.UnpackIntoInterface(&transferBatchEvent, "TransferBatch", vLog.Data)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatal("could not unpack contract Abi")
				}
				//transferBatchEvent
				transferBatchEvent.Operator = common.HexToAddress(vLog.Topics[1].Hex())
				transferBatchEvent.From = common.HexToAddress(vLog.Topics[2].Hex())
				transferBatchEvent.To = common.HexToAddress(vLog.Topics[3].Hex())
			}
		}
	}
}
