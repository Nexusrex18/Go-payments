package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/viper"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)


// Smart contract deployment settings
// const (
// 	gasLimit   = 3000000           // Gas limit for deployment
// 	// gasPrice   = gasPrice       // Gas price (20 Gwei)
// 	chainID    = 17000                 // Ethereum Mainnet (update to match your network)
// 	rpcURL     = "https://holesky.infura.io/v3/6e169b79ad1847e083e71343dfafbf06"   // Replace with your Ethereum node RPC URL
// 	privateKey = "0x84c6835fc30e4a99cbbc0277c1ebec768dc95b1ac731e7ddd2e8ef86fc461825" // Replace with your Ethereum account's private key
// )

// DeployContract deploys a smart contract to the blockchain
func DeployContract() (common.Address, *types.Transaction, error) {
	// Load YAML config
	configFile := "./config/blockchain_config.yaml" // Adjusted path
	viper.SetConfigFile(configFile)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Retrieve values from the YAML configuration
	rpcURL := viper.GetString("blockchain.rpc_url")
	privateKey := viper.GetString("blockchain.private_key")
	gasLimit := viper.GetUint64("blockchain.gas_limit")
	gasPriceGwei := viper.GetInt64("blockchain.gas_price") // Gas price in Gwei
	chainID := viper.GetInt64("blockchain.network_id")

	contractBinPath := viper.GetString("blockchain.contract_bin")

	// Convert Gas Price to Wei
	gasPrice := big.NewInt(gasPriceGwei)
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(1e9))

	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// Load the private key
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Get the account address from the private key
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, nil, fmt.Errorf("failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Read the compiled contract bytecode
	bytecode, err := os.ReadFile(contractBinPath)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to read contract bytecode: %v", err)
	}

	// Create the transaction
	tx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, bytecode)

	// Sign the transaction
	chainIDBig := big.NewInt(chainID)
	signer := types.LatestSignerForChainID(chainIDBig)
	signedTx, err := types.SignTx(tx, signer, privateKeyECDSA)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to wait for transaction mining: %v", err)
	}

	log.Printf("Contract deployed! Address: %s\n", receipt.ContractAddress.Hex())
	return receipt.ContractAddress, signedTx, nil

	// Get the contract address
	// contractAddress := crypto.CreateAddress(fromAddress, nonce)

	// log.Printf("Contract deployed! Address: %s\n", contractAddress.Hex())
	// return contractAddress, signedTx, nil
}


// GetContractABI loads the contract's ABI from a file
func GetContractABI(abiPath string) (abi.ABI, error) {
	abiJSON, err := os.ReadFile(abiPath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI file: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}
	return parsedABI, nil
}

