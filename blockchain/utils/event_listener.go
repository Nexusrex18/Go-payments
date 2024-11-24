package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	// amqp "github.com/rabbitmq/amqp091-go"
	"github.com/Blockchain/config"           // Replace with the correct import path for your config package
	amqp091 "github.com/rabbitmq/amqp091-go" // For RabbitMQ integration
	// "githbub.com/Blockchain/blockchain" // Update with the correct import path for your config package
)

// PaymentEvent represents the structure of the PaymentSent event
type PaymentEvent struct {
	TransactionID string         `json:"transaction_id"`
	Sender        common.Address `json:"sender"`
	Receiver      common.Address `json:"receiver"`
	Amount        *big.Int       `json:"amount"`
	Timestamp     *big.Int       `json:"timestamp"`
	Status        string         `json:"status"` // From payments.proto (e.g., 'SUCCESS', 'PENDING')
}

// EventListener listens for blockchain events and forwards them to the handler
type EventListener struct {
	client            *ethclient.Client
	contractAddress   common.Address
	contractABI       abi.ABI
	eventChannel      chan PaymentEvent
	subscription      ethereum.Subscription
	rabbitMQConn      *amqp091.Connection
	rabbitMQChannel   *amqp091.Channel
	rabbitMQQueueName string
}

// NewEventListener creates a new EventListener instance
func NewEventListener(config *config.BlockchainConfig, eventChan chan PaymentEvent) (*EventListener, error) {
	// Connect to the Ethereum client via WebSocket
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.RabbitMQ.Username,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	client, err := ethclient.Dial(config.Blockchain.WSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum WebSocket: %v", err)
	}

	// Read and parse the contract ABI
	abiData, err := os.ReadFile(config.Blockchain.ContractABI)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract ABI: %v", err)
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	// Connect to RabbitMQ (stub for now, as RabbitMQ isn't in your YAML)
	rabbitMQConn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/") // Update with actual RabbitMQ settings
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	rabbitMQChannel, err := rabbitMQConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}

	// Declare RabbitMQ queue (using a default name for now)
	queueName := "payment_events" // Replace or add this in YAML if RabbitMQ is configured
	_, err = rabbitMQChannel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare RabbitMQ queue: %v", err)
	}

	// Initialize EventListener
	return &EventListener{
		client:            client,
		contractAddress:   common.HexToAddress(config.Blockchain.ContractAddr),
		contractABI:       parsedABI,
		eventChannel:      eventChan,
		rabbitMQConn:      rabbitMQConn,
		rabbitMQChannel:   rabbitMQChannel,
		rabbitMQQueueName: queueName,
	}, nil
}

// StartListening starts the event subscription and listens for logs
func (e *EventListener) StartListening(ctx context.Context) error {
	// Define the query to listen for events from the contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.contractAddress},
	}

	// Subscribe to logs
	logs := make(chan types.Log)
	sub, err := e.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}

	e.subscription = sub
	log.Println("Event listener started. Listening for PaymentSent events...")

	// Listen for logs
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case vLog := <-logs:
				log.Printf("Received log: %+v", vLog)
				// Decode and process the log
				e.processLog(vLog)
			case <-ctx.Done():
				log.Println("Context canceled, stopping event listener.")
				return
			}
		}
	}()

	return nil
}

// processLog decodes the log and forwards the event
func (e *EventListener) processLog(vLog types.Log) {
	// Identify the event based on the topic hash
	eventName, err := e.getEventName(vLog.Topics[0])
	if err != nil {
		log.Printf("Unrecognized event signature: %v", err)
		return
	}

	// Decode the log data
	var event PaymentEvent
	err = e.contractABI.UnpackIntoInterface(&event, eventName, vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack log: %v", err)
		return
	}

	// Decode indexed fields (e.g., sender, receiver) if necessary
	if len(vLog.Topics) > 1 {
		event.Sender = common.HexToAddress(vLog.Topics[1].Hex())
		event.Receiver = common.HexToAddress(vLog.Topics[2].Hex())
	}

	// Enrich the event with status and transaction ID
	event.Status = "PENDING" // Default status; can be updated via gRPC logic
	event.TransactionID = vLog.TxHash.Hex()

	log.Printf("Decoded event: %+v", event)

	// Forward event to the channel
	e.eventChannel <- event

	// Forward event to RabbitMQ (or can be adapted for HTTP callback)
	e.forwardToRabbitMQ(event)
}

// getEventName maps the event topic hash to its name
func (e *EventListener) getEventName(topic common.Hash) (string, error) {
	for name, event := range e.contractABI.Events {
		if topic == event.ID {
			return name, nil
		}
	}
	return "", fmt.Errorf("event not found for topic: %s", topic.Hex())
}

// forwardToRabbitMQ sends the payment event to RabbitMQ for further processing
func (e *EventListener) forwardToRabbitMQ(event PaymentEvent) {
	body := fmt.Sprintf(
		`{"transaction_id": "%s", "sender": "%s", "receiver": "%s", "amount": "%s", "status": "%s", "timestamp": "%s"}`,
		event.TransactionID,
		event.Sender.Hex(),
		event.Receiver.Hex(),
		event.Amount.String(),
		event.Status,
		event.Timestamp.String(),
	)

	err := e.rabbitMQChannel.Publish(
		"",                  // Exchange
		e.rabbitMQQueueName, // Routing key (queue name)
		false,               // Mandatory
		false,               // Immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf("Failed to send event to RabbitMQ: %v", err)
		return
	}

	log.Printf("Sent event to RabbitMQ: %+v", event)
}

// StopListening gracefully stops the event listener
func (e *EventListener) StopListening() {
	if e.subscription != nil {
		e.subscription.Unsubscribe()
		log.Println("Event listener subscription stopped.")
	}

	if e.rabbitMQChannel != nil {
		e.rabbitMQChannel.Close()
	}
	if e.rabbitMQConn != nil {
		e.rabbitMQConn.Close()
	}
}

// func listenForPaymentUpdates(ch *amqp091.Channel, queueName string) {
//     msgs, err := ch.Consume(
//         queueName, // Queue name
//         "",        // Consumer tag
//         false,     // Auto-acknowledge (set to false for manual ack)
//         false,     // Exclusive
//         false,     // No-local
//         false,     // No-wait
//         nil,       // Arguments
//     )
//     if err != nil {
//         log.Fatalf("Failed to register consumer: %v", err)
//     }

//     for msg := range msgs {
//         log.Printf("Received message: %s", msg.Body)

//         // Process the message (e.g., update payment status)
//         err := processMessage(msg.Body)
//         if err != nil {
//             log.Printf("Failed to process message: %v", err)
//             msg.Nack(false, true) // Negative acknowledgment with requeue
//             continue
//         }

//         // Acknowledge successful processing
//         msg.Ack(false)
//         log.Println("Message acknowledged successfully")
//     }
// }
