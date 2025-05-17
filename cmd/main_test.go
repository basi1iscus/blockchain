package main

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/blockchain"
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/sign/sign_ed25519"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"blockchain_demo/pkg/transaction/contract_call"
	"blockchain_demo/pkg/transaction/contract_deploy"
	"blockchain_demo/pkg/transaction/token_transfer"
	"blockchain_demo/pkg/transaction_processor"
	"blockchain_demo/pkg/transaction_processor/coin_transfer_processor"
	"blockchain_demo/pkg/transaction_processor/contract_call_processor"
	"blockchain_demo/pkg/transaction_processor/contract_deploy_processor"
	"blockchain_demo/pkg/transaction_processor/token_transfer_processor"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBlockchainWithDifferentTransactionTypes(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	var addresses = [3]string{"1234567890abcdef1234567890abcdef12345678", "abcdef1234567890abcdef1234567890abcdef12", "2345678901abcdef2345678901abcdef23456789"}
	var signatures = [3]*sign.SignatureKeys{}
	signer := sign_ed25519.Ed25519Signer{}
	for i := 0; i < 3; i++ {
		signatures[i], _ = signer.GenerateKeyPair()
	}
	creator := "ad23947398423423cd234fe34345345323423423"

	// Add BallanceStorage and TransactionProcessor map
	ballanceStorage := ballance_storage.NewMemoryStorage()
	processors := map[transaction.TransactionType]transaction_processor.TransactionProcessor{
		coin_transfer.CoinTransfer:        coin_transfer_processor.NewProcessor(ballanceStorage),
		token_transfer.TokenTransfer:      token_transfer_processor.NewProcessor(ballanceStorage),
		contract_deploy.ContractDeploy:    contract_deploy_processor.NewProcessor(ballanceStorage),
		contract_call.ContractCall:        contract_call_processor.NewProcessor(ballanceStorage),
	}


	// 1 - block
	bc, err := blockchain.NewBlockchain(
		uint64(50000),
		uint64(8),
		creator,
		signer,
		ballanceStorage,
		processors,
	)
	if err != nil {
		t.Fatalf("failed to create blockchain: %v", err)
	}

	// 2 - block

	senderInd := rnd.Intn(len(addresses))
	tx, err := transaction.CreateTransaction(coin_transfer.CoinTransfer, creator, rnd.Int63n(1000), rnd.Int63n(10), map[string]any{
		"recipient": addresses[0],
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	err = tx.AddSing(signer, signatures[senderInd])
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	err = bc.AddTransactionToPool(tx)
	if err != nil {
		t.Fatalf("failed to add transaction to pool: %v", err)
	}

	tx, err = transaction.CreateTransaction(token_transfer.TokenTransfer, creator, rnd.Int63n(1000), rnd.Int63n(10), map[string]any{
		"recipient": addresses[1],
		"token": "2345678901abcdef2345678901abcdef23456789",
		"amount": rnd.Int63n(100),
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	err = tx.AddSing(signer, signatures[senderInd])
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	err = bc.AddTransactionToPool(tx)
	if err != nil {
		t.Fatalf("failed to add transaction to pool: %v", err)
	}

	tx, err = transaction.CreateTransaction(contract_deploy.ContractDeploy, creator, rnd.Int63n(1000), rnd.Int63n(10), map[string]any{
		"code": "2345678901abcdef2345678901abcdef23456789",
		"contractAddress": "abcdef1234567890abcdef1234567890abcdef12",
		"owner": "2345678901abcdef2345678901abcdef23456789",
		"initialSupplay": 1000,		
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	err = tx.AddSing(signer, signatures[senderInd])
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	err = bc.AddTransactionToPool(tx)
	if err != nil {
		t.Fatalf("failed to add transaction to pool: %v", err)
	}

	_, err = bc.MineBlockFromPool(addresses[senderInd])
	if err != nil {
		t.Fatalf("failed to mine block: %v", err)
	}

	err = bc.Verify(4)
	if err != nil {
		t.Errorf("Blockchain verification failed: %v", err)
	}

	// 3 - block
	tx, err = transaction.CreateTransaction(contract_call.ContractCall, creator, rnd.Int63n(1000), rnd.Int63n(10), map[string]any{
		"contractAddress": "abcdef1234567890abcdef1234567890abcdef12",
		"method": "transfer",
		"to": "2345678901abcdef2345678901abcdef23456789",
		"amount": rnd.Int63n(1000),
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	err = tx.AddSing(signer, signatures[senderInd])
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	err = bc.AddTransactionToPool(tx)
	if err != nil {
		t.Fatalf("failed to add transaction to pool: %v", err)
	}

	_, err = bc.MineBlockFromPool(addresses[senderInd])
	if err != nil {
		t.Fatalf("failed to mine block: %v", err)
	}

	err = bc.Verify(4)
	if err != nil {
		t.Errorf("Blockchain verification failed: %v", err)
	}

	fmt.Println(bc)
}
