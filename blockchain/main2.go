// package main2

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	// "time"

// 	"github.com/Blockchain/config" // Update with your actual config import path
// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	// "github.com/"
// 	blockchain "github.com/Blockchain/utils" // Update with your actual blockchain import path

	
// )

// // var blockchainConfig *config.BlockchainConfig
// // var client *ethclient.Client
// // var contractABI abi.ABI
// // var contractAddress common.Address

// func main() {
// 	// Load blockchain configuration
// 	configPath := "./config/blockchain_config.yaml"
// 	blockchainConfig := config.MustLoadConfig(configPath)

// 	// Initialize Ethereum client and contract in the helper
// 	client, contractABI, contractAddress, err := blockchain.InitClient(
// 		blockchainConfig.Blockchain.WSURL,
// 		blockchainConfig.Blockchain.ContractAddr,
// 		blockchainConfig.Blockchain.ContractABI,
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize blockchain client: %v", err)
// 	}
// 	defer client.Close()
// 	log.Println("Blockchain client initialized successfully.")

// 	// Initialize context for graceful shutdown
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Initialize the event listener
// 	eventChannel := make(chan blockchain.PaymentEvent)
// 	eventListener, err := blockchain.NewEventListener(blockchainConfig, eventChannel)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize event listener: %v", err)
// 	}

// 	// Start listening for blockchain events
// 	go func() {
// 		if err := eventListener.StartListening(ctx); err != nil {
// 			log.Printf("Error in event listener: %v", err)
// 		}
// 	}()

// 	log.Println("Listening for blockchain events...")

// 	// Simulated payment logic
// 	go func() {
// 		privateKey := blockchainConfig.Blockchain.PrivateKey
// 		receiver := common.HexToAddress("0x1aF132875b3B0D452c353DB959559C398BEebc40") // Replace with the actual receiver address
// 		amount := big.NewInt(1e16) // Sending 0.1 Ether in Wei

// 		// Parse private key and send payment
// 		senderKey, err := blockchain.ParsePrivateKey(privateKey)
// 		if err != nil {
// 			log.Fatalf("Invalid private key: %v", err)
// 		}

// 		// Send payment using the blockchain utility
// 		signedTx, err := blockchain.SendPayment(
// 			client,
// 			senderKey,
// 			contractAddress,
// 			contractABI,
// 			receiver,
// 			amount,
// 			uint64(blockchainConfig.Blockchain.GasLimit),
// 			big.NewInt(int64(blockchainConfig.Blockchain.NetworkID)),
// 		)
// 		if err != nil {
// 			log.Printf("Failed to send payment: %v", err)
// 			return
// 		}
// 		log.Printf("Payment sent successfully! Transaction hash: %s", signedTx.Hash().Hex())
// 	}()

// 	// Graceful shutdown handling
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

// 	go func() {
// 		<-stop
// 		log.Println("Shutting down service...")
// 		cancel()
// 	}()

// 	// Process blockchain events
// 	go func() {
// 		for event := range eventChannel {
// 			// Example: Process received events
// 			log.Printf("Received blockchain event: %+v", event)
// 		}
// 	}()

// 	// Wait for context cancellation
// 	<-ctx.Done()

// 	// Stop the event listener
// 	eventListener.StopListening()
// 	log.Println("Event listener stopped. Service shut down gracefully.")
// }

// func InitClient(ws_url, contractAddr, contractABIPath string) (*ethclient.Client, abi.ABI, common.Address, error) {
//     client, err := ethclient.Dial(ws_url)
//     if err != nil {
//         return nil, abi.ABI{}, common.Address{}, fmt.Errorf("failed to connect to Ethereum client: %v", err)
//     }

//     abiData, err := os.ReadFile(contractABIPath)
//     if err != nil {
//         return nil, abi.ABI{}, common.Address{}, fmt.Errorf("failed to read contract ABI: %v", err)
//     }

//     parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
//     if err != nil {
//         return nil, abi.ABI{}, common.Address{}, fmt.Errorf("failed to parse contract ABI: %v", err)
//     }

//     return client, parsedABI, common.HexToAddress(contractAddr), nil
// }

// func loadContractABI(abiPath string) abi.ABI {
// 	abiData, err := os.ReadFile(abiPath)
// 	if err != nil {
// 		log.Fatalf("Failed to read contract ABI: %v", err)
// 	}
// 	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
// 	if err != nil {
// 		log.Fatalf("Failed to parse contract ABI: %v", err)
// 	}
// 	return parsedABI
// }

// func listenForEvents(ctx context.Context, client *ethclient.Client, contractABI abi.ABI, contractAddress common.Address, logs chan types.Log) {
// 	// logs = make(chan types.Log)
// 	// defer close(logs)
// 	query := ethereum.FilterQuery{Addresses: []common.Address{contractAddress}}
// 	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
// 	if err != nil {
// 		log.Fatalf("Failed to subscribe to contract logs from main.go: %v", err)
// 	}
// 	defer sub.Unsubscribe()

// 	log.Println("Listening for blockchain events...")

// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			log.Printf("Event subscription error: %v", err)
// 			return
// 		case vLog := <-logs:
// 			// Decode the event log
// 			var event blockchain.PaymentEvent
// 			err := contractABI.UnpackIntoInterface(&event, "PaymentSent", vLog.Data)
// 			if err != nil {
// 				log.Printf("Failed to decode log: %v", err)
// 				continue
// 			}

// 			// Set additional fields (e.g., block number, transaction hash)
// 			// event.BlockNumber = vLog.BlockNumber
// 			// event.TxHash = vLog.TxHash.Hex()

// 			// Send to eventChannel
// 			// eventChannel <- event
// 		case <-ctx.Done():
// 			log.Println("Event listener stopped.")
// 			return
// 		}
// 	}
// }



// // func processLog(vLog types.Log, contractABI abi.ABI) {
// // 	eventName, err := GetEventName(contractABI, vLog.Topics[0])
// // 	if err != nil {
// // 		log.Printf("Unknown event: %v", err)
// // 		return
// // 	}
// // 	log.Printf("Event detected: %s - Data: %x", eventName, vLog.Data)
// // }

// // func GetEventName(parsedABI abi.ABI, topic common.Hash) (string, error) {
// //     for name, event := range parsedABI.Events {
// //         if topic == event.ID {
// //             return name, nil
// //         }
// //     }
// //     return "", fmt.Errorf("event not found for topic: %s", topic.Hex())
// // }

// // func SendPayment(receiver common.Address, amount *big.Int) error {
// // 	// Prepare transaction data
// // 	contractABIPath := blockchainConfig.Blockchain.ContractABI
// // 	contractABIContent, err := os.ReadFile(contractABIPath)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to read contract ABI file: %v", err)
// // 	}

// // 	parsedABI, err := abi.JSON(strings.NewReader(string(contractABIContent)))
// // 	if err != nil {
// // 		return fmt.Errorf("failed to parse contract ABI: %v", err)
// // 	}

// // 	// _, err := parsedABI.Pack("sendPayment", receiver)
// // 	// if err != nil {
// // 	// 	return fmt.Errorf("failed to pack data: %v", err)
// // 	// }
// // 	data, err := parsedABI.Pack("sendPayment", receiver)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to pack data: %v", err)
// // 	}

// // 	// Estimate gas
// // 	gasLimit := blockchainConfig.Blockchain.GasLimit
// // 	gasPrice := blockchainConfig.Blockchain.GasPrice

// // 	log.Printf("Estimated gas: %d, gas price: %d", gasLimit, gasPrice)
// // 	// Add signing and transaction sending logic here
// // 	return nil
// // } 
// // hello
