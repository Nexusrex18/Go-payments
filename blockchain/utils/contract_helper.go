package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"
	"log"
	// "time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum"
)

var (
	client          *ethclient.Client
	contractAddress common.Address
	parsedABI       abi.ABI
)

// InitClient initializes the Ethereum client, parses the ABI, and sets the contract address.
func InitClient(rpcURL, contractAddr, contractABIPath string) (*ethclient.Client, abi.ABI, common.Address, error) {
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, abi.ABI{}, common.Address{}, fmt.Errorf("failed to connect to Ethereum client: %v", err)
    }

    abiData, err := os.ReadFile(contractABIPath)
    if err != nil {
        return nil, abi.ABI{}, common.Address{}, fmt.Errorf("failed to read contract ABI: %v", err)
    }

    parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
    if err != nil {
        return nil, abi.ABI{}, common.Address{}, fmt.Errorf("failed to parse contract ABI: %v", err)
    }

    return client, parsedABI, common.HexToAddress(contractAddr), nil
}

// SendPayment sends ETH from the sender to the receiver using the smart contract.
func SendPayment(
	client *ethclient.Client,
	senderKey *ecdsa.PrivateKey,
	contractAddress common.Address,
	parsedABI abi.ABI,
	receiverAddress common.Address,
	amount *big.Int,
	gasLimit uint64,
	chainID *big.Int,
) (*types.Transaction, error) {
	// Get sender address
	senderAddress := crypto.PubkeyToAddress(senderKey.PublicKey)
	fmt.Printf("Sending payment to receiver: %s\n", receiverAddress.Hex())

	// Get nonce for the sender
	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get gas price from the network
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// Pack contract function call with parameters
	data, err := parsedABI.Pack("sendPayment", receiverAddress, amount) // Adjust function signature if needed
	if err != nil {
		return nil, fmt.Errorf("failed to pack transaction data: %v", err)
	}

	// Estimate gas if not provided
	if gasLimit == 0 {
		callMsg := ethereum.CallMsg{
			From:     senderAddress,
			To:       &contractAddress,
			GasPrice: gasPrice,
			Value:    big.NewInt(0), // Ether value sent (0 unless required)
			Data:     data,
		}
		gasLimit, err = client.EstimateGas(context.Background(), callMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %v", err)
		}
	}

	// Create the transaction (value set to 0 unless sending Ether directly)
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// Sign the transaction
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, senderKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction to the Ethereum network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	// Return the signed transaction
	return signedTx, nil
}



// GetPaymentDetails retrieves the details of a payment using its ID.
func GetPaymentDetails(paymentID uint64) (common.Address, common.Address, *big.Int, uint64, error) {
	// Pack function call data
	callData, err := parsedABI.Pack("getPaymentDetails", paymentID)
	if err != nil {
		return common.Address{}, common.Address{}, nil, 0, fmt.Errorf("failed to pack call data: %v", err)
	}

	// Call the smart contract
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return common.Address{}, common.Address{}, nil, 0, fmt.Errorf("failed to call contract: %v", err)
	}

	// Decode the result
	var decoded []interface{}
	err = parsedABI.UnpackIntoInterface(&decoded, "getPaymentDetails", result)
	if err != nil {
		return common.Address{}, common.Address{}, nil, 0, fmt.Errorf("failed to unpack contract response: %v", err)
	}

	sender := decoded[0].(common.Address)
	receiver := decoded[1].(common.Address)
	amount := decoded[2].(*big.Int)
	timestamp := decoded[3].(uint64)

	return sender, receiver, amount, timestamp, nil
}

// ContractBalance retrieves the balance of the contract.
func ContractBalance() (*big.Int, error) {
	// Pack function call data
	callData, err := parsedABI.Pack("contractBalance")
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %v", err)
	}

	// Call the smart contract
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	// Decode the result
	var balance *big.Int
	err = parsedABI.UnpackIntoInterface(&balance, "contractBalance", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack contract response: %v", err)
	}

	return balance, nil
}

// ParsePrivateKey parses a hex-encoded private key string.
func ParsePrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(hexKey, "0x"))
}

// ListenForEvents subscribes to contract events and handles them.
func ListenForEvents(ctx context.Context, logs chan types.Log) {
	query := ethereum.FilterQuery{Addresses: []common.Address{contractAddress}}
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to contract logs from contracthelper: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("Listening for blockchain events...")

	for {
		select {
		case vLog := <-logs:
			processLog(vLog)
		case err := <-sub.Err():
			log.Printf("Subscription error: %v", err)
		case <-ctx.Done():
			log.Println("Context canceled, stopping listener.")
			return
		}
	}
}

// processLog processes individual logs received from the blockchain.
func processLog(vLog types.Log) {
	eventName, err := getEventName(vLog.Topics[0])
	if err != nil {
		log.Printf("Unknown event: %v", err)
		return
	}

	log.Printf("Event detected: %s - Data: %x", eventName, vLog.Data)
}

func getEventName(topic common.Hash) (string, error) {
	for name, event := range parsedABI.Events {
		if topic == event.ID {
			return name, nil
		}
	}
	return "", fmt.Errorf("event not found for topic: %s", topic.Hex())
}


// func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
//     for {
//         receipt, err := client.TransactionReceipt(context.Background(), txHash)
//         if receipt != nil {
//             return receipt, nil
//         }
//         if err != nil {
//             log.Printf("Error getting receipt: %v", err)
//         }
//         time.Sleep(1 * time.Second) // Polling interval
//     }
// }