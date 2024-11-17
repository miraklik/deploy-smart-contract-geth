package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx := context.Background()
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/dfa0335a8d2b4364bd669159aa3dc734")
	if err != nil {
		log.Fatal("Unable to connect to node", err)
	}

	log.Println("Connected to node")

	privateKeyHex, exists := os.LookupEnv("PRIVATE_KEY")
	if !exists {
		log.Fatal("PRIVATE_KEY environment variable is not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("Unable to get private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Unable to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal("Unable to get nonce", err)
	}

	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("Unable to get gas price", err)
	}

	toAddress := common.HexToAddress(GenerateRandomWalletAddress())
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal("Unable to get network ID", err)
	}

	signetTx, err := types.SignTx(tx, types.EIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("Unable to sign transaction", err)
	}

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatal("Unable to send transaction", err)
	}

	log.Printf("Transaction sent: %s", signetTx.Hash().Hex())
}

func GenerateRandomWalletAddress() string {
	const walletLength = 20 // это в байтах (20 байт = 40 символа)
	bytes := make([]byte, walletLength)

	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("Failed to generate random bytes: %v", err))
	}
	r
	return "0x" + hex.EncodeToString(bytes)
}
