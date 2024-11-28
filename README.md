# Scalar Core

Scalar Core is a blockchain interoperability network that enables seamless asset transfers between Bitcoin and other chains. Built on top of the Cosmos SDK, Scalar introduces specialized modules for managing cross-chain assets and protocols.

1. First version of the code is based on the axelar codebase
2. Second version of the code is brand new and based on the cosmos sdk and reuses some axelar packages

## Overview

Scalar serves as a decentralized bridge between Bitcoin and various EVM chains and non-EVM chains in the future, providing secure and efficient cross-chain communication. The network utilizes a validator set to manage assets and sign transactions across different chains.

## Key Features

- Bitcoin to EVM bridge functionality
- Multi-chain support for EVM-compatible networks
- Secure multi-signature schemes (Schnorr for BTC, ECDSA for EVM)
- Protocol-level asset management
- Covenant-based custody system

## Core Modules

### 🌟 Covenant Module

Manages the secure custody of user assets across Bitcoin and EVM chains:

- Custodian management system
- Transaction signing for BTC unstaking
- EVM transaction signing for staking operations
- Asset security and management

### 🌟 Protocol Module

Acts as the service layer for managing staking operations:

- ERC20 token management for BTC staking
- Protocol information storage on Scalar network
- Cross-chain communication coordination

### 🌟 EVM ERC20 Tokens Module

Handles the deployment and management of ERC20 tokens across EVM chains:

- Multiple token deployment capability per chain
- Configuration for:
  - Bitcoin source network (testnet4, mainnet, regtest)
  - Protocol identity (protocol pubkey)
  - Covenant management settings

### 🌟 Bitcoin Module

Provides Bitcoin network integration:

- Transaction verification
- Transaction validation
- Confirmation management
- Bitcoin network state tracking

### 🌟 Multisig Module

Implements secure multi-signature schemes:

- EVM: ECDSA multisig with weighted signatures
- BTC: Schnorr multisig with Taproot support
- Validator coordination for transaction signing

## Prerequisites

- Go 1.23+
- Rust 1.82+
- Docker

## Building

```bash
make build
```

## Running

```bash
make start
```

## Docker

### Building Docker Images

```bash
make docker-image
```

#### Running Docker Container

```bash
make docker-run
```

## Development Setup

1. Clone the repository:

```bash
git clone https://github.com/scalarorg/xchains-core.git
cd xchains-core
```

2. Install dependencies:

```bash
make prereqs
```

3. Generate protocol buffers:

```bash
make proto-gen
```

4. Build Docker Image:

```bash
make docker-image
```

5. Run Docker Container:

```bash
make docker-run
```

## Configuration

The network can be configured through environment variables or a config file. Key configuration options:

- `NODE_MONIKER`: The Scalar node's moniker
- `PEERS_FILE`: File with peer list for network connection
- `CONFIG_PATH`: Path to configuration file
- `PRESTART_SCRIPT`: Pre-launch script path

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

## License

[License details to be added]

## Security

For security concerns, please email [security contact to be added].
