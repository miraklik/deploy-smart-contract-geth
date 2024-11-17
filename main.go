package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/dfa0335a8d2b4364bd669159aa3dc734")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyHex, exists := os.LookupEnv("PRIVATE_KEY")
	if !exists {
		log.Fatal("PRIVATE_KEY environment variable is not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	byteCodeCheck, exists := os.LookupEnv("BYTECODE_CHECK")
	if !exists {
		log.Fatal("BYTECODE_CHECK environment variable is not set")
	}

	bytecode := common.FromHex(byteCodeCheck)

	address, tx, _, err := bind.DeployContract(auth, bind.ContractBackend(client), nil, bytecode, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Contract deployed! Address: %s, Transaction: %s", address.Hex(), tx.Hash().Hex())
}
