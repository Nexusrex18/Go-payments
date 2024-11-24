package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// BlockchainConfig represents the structure of the blockchain_config.yaml file
type BlockchainConfig struct {
	Blockchain struct {
		RPCURL       string `yaml:"rpc_url"`       // RPC URL for Ethereum (e.g., Infura/Alchemy or local node)
		WSURL        string `yaml:"ws_url"`        // WebSocket URL for event listening
		ContractABI  string `yaml:"contract_abi"`  // Path to the contract ABI file
		ContractBin  string `yaml:"contract_bin"`  // Path to the contract bytecode file
		ContractAddr string `yaml:"contract_addr"` // Deployed contract address
		PrivateKey   string `yaml:"private_key"`   // Private key for signing transactions
		GasLimit     uint64 `yaml:"gas_limit"`     // Gas limit for transactions
		GasPrice     int64  `yaml:"gas_price"`     // Gas price in Gwei
		NetworkID    int64  `yaml:"network_id"`    // Ethereum network ID (e.g., 1 for mainnet, 5 for Goerli)
	} `yaml:"blockchain"`
	RabbitMQ struct {
		Host     string `yaml:"host"`     // RabbitMQ host
		Port     int    `yaml:"port"`     // RabbitMQ port
		Username string `yaml:"username"` // RabbitMQ username
		Password string `yaml:"password"` // RabbitMQ password
		Exchange string `yaml:"exchange"` // RabbitMQ exchange name
		Queue    string `yaml:"queue"`    // RabbitMQ queue name
	} `yaml:"rabbitmq"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(filePath string) (*BlockchainConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config BlockchainConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	return &config, nil
}

// MustLoadConfig loads the configuration and exits on failure
func MustLoadConfig(filePath string) *BlockchainConfig {
	config, err := LoadConfig(filePath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Println("Configuration loaded successfully.")
	return config
}
