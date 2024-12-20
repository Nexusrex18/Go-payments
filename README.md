
# Go-Payments

## Overview
**Blockchain Odyssey** is a cross-border payment service that leverages blockchain technology to provide secure, transparent, and scalable transactions. The project is built using **Go**, integrating **Ethereum** or **Solana** smart contracts, along with GoFr microservices for transaction management, user wallet management, and blockchain integration.

## Key Features:
- **Secure & Transparent Payments**: Built with blockchain technology to ensure secure and transparent transactions.
- **Currency Conversion**: Handles cross-border payments and automatic conversion of currencies.
- **Real-Time Transaction Tracking**: Users can track transactions in real time for enhanced transparency.
- **Scalable & Modular Architecture**: Modular GoFr microservices ensure the system is scalable and can easily handle a large volume of transactions.

## Project Structure

```
C:.
|   README.md
|
+---blockchain
|   |   docker-compose.yaml
|   |   go.mod
|   |   go.sum
|   |   main.go
|   |   main2.go
|   |
|   +---config
|   |       blockchain_config.yaml
|   |       config.go
|   |
|   +---contracts
|   |       compile.sh
|   |       Payment.abi
|   |       Payment.bin
|   |       Payments.sol
|   |
|   +---deployment
|   |       deploy.go
|   |       migration.json
|   |
|   \---utils
|       |   contract_helper.go
|       |   event_listener.go
|       |
|       \---abi
+---payment-service
|   |   .gitignore
|   |   docker-compose.yml
|   |   Dockerfile
|   |   go.mod
|   |   go.sum
|   |   main.go
|   |
|   +---cmd
|   |   \---payments
|   |           main.go
|   |
|   \---internal
|       |   .DS_Store
|       |
|       +---api
|       |   +---error
|       |   |       errors.go
|       |   |
|       |   +---grpc
|       |   |   |   payments.proto
|       |   |   |   payments_handler.go
|       |   |   |
|       |   |   \---main
|       |   |           server.go
|       |   |
|       |   \---middleware
|       |           auth.go
|       |
|       +---config
|       |       config.go
|       |
|       +---db
|       |       db.go
|       |
|       +---proto
|       |   |   payments.pb.go
|       |   |   payments_grpc.pb.go
|       |   |
|       |   \---grpc
|       |           payments.pb.go
|       |           payments_grpc.pb.go
|       |
|       \---rabbitmq
|               connection.go
|
\---user-authentication
    |   go.mod
    |   go.sum
    |   main.go
    |
    +---db
    |       db.go
    |
    \---jwt-tokenization
            tokens.go
```

## Project Breakdown

### Blockchain Service
This part handles the integration of blockchain for secure transactions. The smart contract code, blockchain deployment scripts, and utilities for interacting with the blockchain are housed here.

- **Main Files**:
  - `main.go`: Entry point for blockchain interaction.
  - `contracts/`: Contains Ethereum/Solana smart contract code.
  - `utils/`: Utilities for contract interaction and event listening.

### Payment Service
This service handles the backend functionality for payment processing, including transaction handling and communication via gRPC APIs.

- **Main Files**:
  - `main.go`: Entry point for payment service API.
  - `api/`: Contains the core API logic, including error handling and gRPC services.
  - `rabbitmq/`: Manages communication between services via RabbitMQ.

### User Authentication
This module handles user authentication, including JWT token generation and validation.

- **Main Files**:
  - `tokens.go`: Handles JWT token generation and validation.
  - `db/`: Contains database-related code for user management.

## Getting Started

### Prerequisites
- **Go** installed (version 1.16 or above)
- **Docker** and **Docker Compose** for managing containers
- **Solidity** for smart contract development (if using Ethereum)

### Setting Up

1. Clone the repository:
   ```bash
   git clone https://github.com/Nexusrex18/go-payments.git
   cd go-payments
   cd payment-service
   ```

2. Build and run the Docker containers:
   ```bash
   docker-compose up --build
   ```

3. Run the payment service:
   ```bash
   go run main.go
   ```

4. Deploy smart contracts using:
   ```bash
   ./contracts/compile.sh
   ```

### Testing the Service

- Once the services are up, you can interact with the payment service via the exposed gRPC API.
- For user authentication, use the `/auth` endpoint to register and log in.

## Contributing
Feel free to fork the repository, create a branch, and submit a pull request for any enhancements, bug fixes, or new features.
