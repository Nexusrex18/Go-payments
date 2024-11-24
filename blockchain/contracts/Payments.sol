// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Payment {
    address public owner;

    // Event to log payment details
    event PaymentSent(address indexed sender, address indexed receiver, uint256 amount, uint256 timestamp);

    // Struct to store payment details
    struct PaymentDetail {
        address sender;
        address receiver;
        uint256 amount;
        uint256 timestamp;
    }

    // Mapping to store payments by ID
    mapping(uint256 => PaymentDetail) public payments;

    // Counter for unique payment IDs
    uint256 public paymentCount;

    // Modifier to restrict functions to the owner
    modifier onlyOwner() {
        require(msg.sender == owner, "Action restricted to the contract owner");
        _;
    }

    // Modifier to validate receiver address
    modifier validAddress(address _receiver) {
        require(_receiver != address(0), "Receiver address cannot be zero");
        require(_receiver != address(this), "Receiver cannot be the contract itself");
        _;
    }

    // Constructor to set the deployer as the contract owner
    constructor() {
        owner = msg.sender;
    }

    // Function to send payment to a receiver
    function sendPayment(address payable _receiver)
        external
        payable
        validAddress(_receiver)
        returns (uint256)
    {
        require(msg.value > 0, "Payment amount must be greater than zero");

        // Increment payment counter
        paymentCount++;

        // Save payment details in storage
        payments[paymentCount] = PaymentDetail({
            sender: msg.sender,
            receiver: _receiver,
            amount: msg.value,
            timestamp: block.timestamp
        });

        // Emit payment event
        emit PaymentSent(msg.sender, _receiver, msg.value, block.timestamp);

        // Transfer funds to receiver
        (bool success, ) = _receiver.call{value: msg.value}("");
        require(success, "Payment transfer failed");

        return paymentCount;
    }

    // Function to fetch payment details by ID
    function getPaymentDetails(uint256 paymentId)
        external
        view
        returns (address sender, address receiver, uint256 amount, uint256 timestamp)
    {
        require(paymentId > 0 && paymentId <= paymentCount, "Invalid payment ID");

        PaymentDetail memory payment = payments[paymentId];
        return (payment.sender, payment.receiver, payment.amount, payment.timestamp);
    }

    // Function to view contract's balance
    function contractBalance() external view onlyOwner returns (uint256) {
        return address(this).balance;
    }

    // Owner-only function to withdraw funds
    function withdraw(uint256 _amount) external onlyOwner {
        require(_amount > 0, "Withdrawal amount must be greater than zero");
        require(_amount <= address(this).balance, "Insufficient contract balance");

        payable(owner).transfer(_amount);
    }
}
