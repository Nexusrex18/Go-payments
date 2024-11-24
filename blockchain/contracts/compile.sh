#!/bin/bash

# Path to your Solidity contract
CONTRACT_PATH="./contracts/Payments.sol"

# Output directory for the compiled contract artifacts
OUTPUT_DIR="./contracts"

# Check if solc-windows is installed and accessible
if ! command -v solc-windows &> /dev/null
then
    echo "Error: Solidity compiler 'solc-windows' not found. Please install it and ensure it's in your PATH."
    exit 1
fi

# Compile the Solidity contract
echo "Compiling contract: $CONTRACT_PATH"

solc-windows --optimize --abi --bin --overwrite --output-dir $OUTPUT_DIR $CONTRACT_PATH 2>&1 | tee compile.log

# Verify that ABI and BIN files were generated successfully
ABI_FILE="$OUTPUT_DIR/Payments.abi"
BIN_FILE="$OUTPUT_DIR/Payments.bin"

if [ -f "$ABI_FILE" ] && [ -f "$BIN_FILE" ]; then
    echo "Contract compiled successfully."
    echo "ABI file: $ABI_FILE"
    echo "BIN file: $BIN_FILE"
else
    echo "Error: Contract compilation failed. Check compile.log for details."
    exit 1
fi
