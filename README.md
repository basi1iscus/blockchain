# Blockchain Demo

This is a simple blockchain implementation in Go, demonstrating basic blockchain, block, transaction, and signature logic.

## Project Structure

- `cmd/main.go` — Example entry point for running the blockchain demo
- `pkg/blockchain/` — Blockchain logic
- `pkg/block/` — Block structure and mining
- `pkg/transaction/` — Transaction structure and signing
  - `coin_transfer/` — Coin transfer transaction type
  - `contract_call/` — Smart contract call transaction type
  - `contract_deploy/` — Smart contract deployment transaction type
  - `token_transfer/` — Token transfer transaction type
- `pkg/sign/` — Signature generation and verification
- `pkg/utils/` — Utility functions

## Requirements

- Go 1.24.2 or later

## Running Tests

To run all tests in the project, execute:

```
go test ./...
```

## Running the Demo

To run the example main program:

```
go run ./cmd/main.go
```

## Features

- ECDSA key generation and signatures
- Transaction creation, signing, and verification
- Dynamic transaction type registry (реестр типов транзакций)
- Serialization and deserialization of transactions
- Block mining with adjustable difficulty
- Blockchain validation
- Support for multiple transaction types (coin transfer, contract call, contract deploy, token transfer)

## Adding New Transaction Types

1. Create a new subfolder in `pkg/transaction/` for your type.
2. Implement the transaction struct and logic.
3. Register the type in the package's `init()`.
4. Import the package (with `_` if needed) where you use dynamic creation.

## License

MIT License
