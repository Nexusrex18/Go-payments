package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var url = "https://holesky.infura.io/v3/6e169b79ad1847e083e71343dfafbf06"

func main() {
	// Connect to Ethereum client
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// User input for fiat-to-crypto conversion
	var fiatAmount float64
	var senderFiat, receiverFiat string

	fmt.Print("Enter amount in sender's fiat currency: ")
	fmt.Scanln(&fiatAmount)
	fmt.Print("Enter sender's fiat currency (e.g., usd): ")
	fmt.Scanln(&senderFiat)
	fmt.Print("Enter receiver's fiat currency (e.g., eur): ")
	fmt.Scanln(&receiverFiat)

	// Get conversion rate from fiat to Ethereum
	ethRate, err := getConversionRate(senderFiat, "ethereum")
	if err != nil {
		log.Fatalf("Failed to fetch Ethereum conversion rate: %v", err)
	}

	// Calculate equivalent Ether in Wei
	amountInEther := big.NewFloat(fiatAmount / ethRate)
	amountInWei := new(big.Int)
	amountInEther.Mul(amountInEther, big.NewFloat(1e18)).Int(amountInWei)

	// Wallet details
	sender := "6c50E2d7ddB983451bCab2438D3Ed03E5F01B2cE"
	receiver := "1aF132875b3B0D452c353DB959559C398BEebc40"
	privateKey := "84c6835fc30e4a99cbbc0277c1ebec768dc95b1ac731e7ddd2e8ef86fc461825"

	// Check sender and receiver balances
	checkBalances(client, sender, receiver)

	// Send transaction
	txHash, err := sendTransaction(client, sender, receiver, amountInWei, privateKey)
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}
	fmt.Printf("Transaction successful! Hash: %s\n", txHash)

	// Check balances again
	checkBalances(client, sender, receiver)

	// Calculate amount received in receiver's fiat currency
	receiverRate, err := getConversionRate("ethereum", receiverFiat)
	if err != nil {
		log.Fatalf("Failed to fetch receiver's fiat conversion rate: %v", err)
	}
	amountReceived := new(big.Float).Mul(new(big.Float).SetInt(amountInWei), big.NewFloat(receiverRate/1e18))
	fmt.Printf("Receiver will receive approximately %.2f %s\n", amountReceived, receiverFiat)
}

func checkBalances(client *ethclient.Client, sender, receiver string) {
	senderBalance, err := getWalletBalance(client, sender)
	if err != nil {
		log.Fatalf("Failed to fetch sender balance: %v", err)
	}
	fmt.Printf("Sender balance: %s Ether\n", weiToEther(senderBalance))

	receiverBalance, err := getWalletBalance(client, receiver)
	if err != nil {
		log.Fatalf("Failed to fetch receiver balance: %v", err)
	}
	fmt.Printf("Receiver balance: %s Ether\n", weiToEther(receiverBalance))
}

func getWalletBalance(client *ethclient.Client, address string) (*big.Int, error) {
	walletAddress := common.HexToAddress(address)
	return client.BalanceAt(context.Background(), walletAddress, nil)
}

func weiToEther(wei *big.Int) string {
	ether := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return ether.Text('f', 18)
}

func sendTransaction(client *ethclient.Client, sender, receiver string, amount *big.Int, privateKeyHex string) (string, error) {
	senderAddress := common.HexToAddress(sender)
	receiverAddress := common.HexToAddress(receiver)

	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		return "", fmt.Errorf("failed to fetch nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %v", err)
	}

	tx := types.NewTransaction(nonce, receiverAddress, amount, 21000, gasPrice, nil)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

func getConversionRate(fromCurrency, toCurrency string) (float64, error) {
	apiURL := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", fromCurrency, toCurrency)
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result[fromCurrency][toCurrency], nil
}
