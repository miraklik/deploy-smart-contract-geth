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
		panic(err)
	}
	defer client.Close()

	privateKeyHex, exists := os.LookupEnv("PRIVATE_KEY")
	if !exists {
		log.Fatal("PRIVATE_KEY environment variable is not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal("Failed to get pending nonce:", err)
	}

	value := big.NewInt(1000000000000000000) // 1 ETH in wei
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("Failed to suggest gas price:", err)
	}

	toAddress := common.HexToAddress(GenerateRandomWalletAddress())
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction sent: %s", signedTx.Hash().Hex())
}

func GenerateRandomWalletAddress() string {
	const walletLength = 20 // это в байтах (20 байт = 40 символа)
	bytes := make([]byte, walletLength)

	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("Failed to generate random bytes: %v", err))
	}

	return "0x" + hex.EncodeToString(bytes)
}
