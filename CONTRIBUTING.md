
# Contributing to Go-Payments

Thank you for your interest in contributing to **Go-Payments**! We welcome contributions to improve this project, whether it's bug fixes, feature enhancements, or documentation updates.

## Table of Contents
- [Getting Started](#getting-started)
- [Reporting Issues](#reporting-issues)
- [Creating a Pull Request](#creating-a-pull-request)
- [Coding Guidelines](#coding-guidelines)
- [Code of Conduct](#code-of-conduct)

---

## Getting Started

To get started with contributing, follow these steps:

1. **Fork the repository**:  
   Visit the project repository and click the "Fork" button to create your own copy.

2. **Clone your fork**:  
   ```bash
   git clone https://github.com/Nexusrex18/go-payments.git
   cd go-payments
   ```

3. **Set up the environment**:  
   Make sure you have the following prerequisites installed:
   - Go (1.16 or above)
   - Docker and Docker Compose
   - Solidity (for smart contract development)
   - PostgreSQL (for database setup)

4. **Create a new branch**:  
   Create a branch for your contribution:  
   ```bash
   git checkout -b feature/your-feature-name
   ```

5. **Run the project locally**:  
   Use Docker Compose to set up the services:  
   ```bash
   docker-compose up --build
   ```

6. **Write and test your code**:  
   Follow the coding guidelines below and ensure your changes are thoroughly tested.

---

## Reporting Issues

If you encounter bugs or have feature requests, please report them by creating an issue in the repository. Make sure to include:
- A clear title and description.
- Steps to reproduce the issue.
- Environment details (Go version, Docker version, etc.).

---

## Creating a Pull Request

Once your changes are ready, follow these steps to submit your contribution:

1. **Ensure your branch is up to date**:  
   Sync your fork with the main repository:  
   ```bash
   git fetch upstream
   git merge upstream/main
   ```

2. **Run tests**:  
   Ensure all tests pass and your changes do not break existing functionality.

3. **Commit your changes**:  
   Use a descriptive commit message:  
   ```bash
   git commit -m "Add feature: detailed payment status tracking"
   ```

4. **Push to your branch**:  
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Submit a pull request**:  
   - Go to the original repository.
   - Click "New Pull Request."
   - Provide a detailed description of your changes and the problem they solve.

---

## Coding Guidelines

To maintain code quality, please follow these guidelines:

1. **Code Style**:
   - Use **gofmt** to format Go code.
   - Follow idiomatic Go practices.
   - Ensure Solidity contracts are clean and well-commented.

2. **Documentation**:
   - Comment all exported functions and methods.
   - Update `README.md` or related documentation if your changes affect functionality.

3. **Testing**:
   - Write unit tests for new features and bug fixes.
   - Ensure existing tests pass:  
     ```bash
     go test ./...
     ```

4. **Modular Code**:
   - Follow the project’s modular architecture:
     - **Blockchain** for smart contracts and blockchain interactions.
     - **Payment-Service** for backend APIs and RabbitMQ communication.
     - **User-Authentication** for user management and JWT handling.

5. **Dockerized Setup**:
   - Ensure new features are compatible with the Docker environment.
   - Update `docker-compose.yml` if necessary.

---

## Code of Conduct

By participating in this project, you agree to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md), which promotes a respectful and inclusive environment for all contributors.

---

Thank you for your contributions and helping to make **Go-Payments** better!
