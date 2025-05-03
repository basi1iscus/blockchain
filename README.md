# Blockchain Demo

This is a simple blockchain implementation in Go, demonstrating basic blockchain, block, transaction, and signature logic.

## Project Structure

- `cmd/main.go` — Example entry point for running the blockchain demo
- `pkg/blockchain/` — Blockchain logic
- `pkg/block/` — Block structure and mining
- `pkg/transaction/` — Transaction structure and signing
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
- Block mining with adjustable difficulty
- Blockchain validation

## License

MIT License
