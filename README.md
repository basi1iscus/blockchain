# Blockchain Demo

This is a simple blockchain implementation in Go, demonstrating basic blockchain, block, transaction, and signature logic.

## Project Structure

- `cmd/main.go` — Example entry point for running the blockchain demo
- `cmd/main_test.go` — Integration test for blockchain with all transaction types
- `pkg/blockchain/` — Blockchain logic
- `pkg/block/` — Block structure and mining
- `pkg/merkle/` — Merkle tree and root calculation
- `pkg/transaction/` — Transaction structure and signing
  - `coin_transfer/` — Coin transfer transaction type
  - `contract_call/` — Smart contract call transaction type
  - `contract_deploy/` — Smart contract deployment transaction type
  - `token_transfer/` — Token transfer transaction type
- `pkg/transaction_processor/` — Transaction verification and processors for each type
- `pkg/sign/` — Signature generation and verification
- `pkg/utils/` — Utility functions
  - `stack/` — Generic stack implementation (tested, used by script VM)
  - `queue/` — Generic queue implementation (tested, used by script VM for opcode precompilation)
- `pkg/wallet/` — Wallet creation, address validation, and tests
- `pkg/ballance_storage/` — In-memory balance storage and tests
- `pkg/script_vm/` — Bitcoin-like Script VM (stack-based, supports custom opcodes, queue-based precompilation, and signature/hash operations)

## Requirements

- Go 1.24.2 or later

## Running Tests

To run all tests in the project, execute:

```
go test ./...
```

This will run unit tests for all packages, including:
- Blockchain logic
- Transaction types and processors
- Wallet creation and address validation
- Balance storage
- Integration test for all transaction types in `cmd/main_test.go`

## Running the Demo

To run the example main program:

```
go run ./cmd/main.go
```

## Features

- ECDSA and Ed25519 key generation and signatures
- Transaction creation, signing, and verification
- Dynamic transaction type registry (реестр типов транзакций)
- Serialization and deserialization of transactions
- Block mining with adjustable difficulty
- Blockchain validation
- Merkle tree and root calculation for block transactions
- Support for multiple transaction types (coin transfer, contract call, contract deploy, token transfer)
- Wallet creation and address validation
- In-memory balance storage

## Adding New Transaction Types

1. Create a new subfolder in `pkg/transaction/` for your type.
2. Implement the transaction struct and logic.
3. Register the type in the package's `init()`.
4. Implement a processor in `pkg/transaction_processor/` and add it to the processor map.
5. Import the package (with `_` if needed) where you use dynamic creation.

## Blockchain Construction Note

When constructing a blockchain, you must now provide:
- A `ballance_storage.BallanceStorage` implementation (e.g., `NewMemoryStorage()`)
- A map of transaction processors for all supported types (see `cmd/main_test.go` for an example)

## Script VM, Stack, and Queue

### Script VM (`pkg/script_vm/`)
- Implements a Bitcoin-like stack-based virtual machine for transaction scripts.
- Supports opcode precompilation using a queue for efficient execution.
- Handles standard stack operations, signature/hash opcodes, and custom logic.
- Extensible for new opcodes and script types.
- Used for validating P2PKH, multisig, and custom scripts.

### Stack (`pkg/utils/stack/`)
- Generic, type-safe stack implementation in Go.
- Used by the script VM for all stack operations.
- Fully unit tested.

### Queue (`pkg/utils/queue/`)
- Generic, type-safe queue implementation in Go.
- Used by the script VM for opcode precompilation and execution order.
- Provides iterator support for idiomatic Go iteration.
- Fully unit tested.

## API Endpoints

The demo provides a simple HTTP API using [Gin](https://github.com/gin-gonic/gin):

### POST `/api/wallet`
- **Description:** Create a new wallet (Ed25519 keypair, address, public key hash).
- **Response:**
  - `address`: Wallet address (Base58)
  - `public_key`: Public key (hex)
  - `public_key_hash`: Public key hash (hex)
  - `private_key`: Private key (hex)

### POST `/api/run_sript`
- **Description:** Run a script using the Script VM. Accepts a scriptSig, scriptPubKey, and signed data, executes the script, and returns the parsed script code.
- **Request JSON:**
  - `script_sig`: ScriptSig as hex string
  - `script_pub_key`: ScriptPubKey as hex string
  - `signed_data`: Data to be verified (hex string)
- **Response:**
  - Returns the parsed script code as a string (or error message if execution fails).

### GET `/ping`
- **Description:** Health check endpoint. Returns `{ "message": "pong" }`.

---

You can use tools like [Postman](https://www.postman.com/) or `curl` to interact with these endpoints.

## License

MIT License
