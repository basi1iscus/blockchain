package blockchain

import (
	"blockchain_demo/pkg/block"
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/sign/sign_ecdsa"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"encoding/hex"
	"testing"
)

func generateTestKeys(t *testing.T, signer sign.Signer) *sign.SignatureKeys {
	signature, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	return signature
}

func randomAddress() string {
	return hex.EncodeToString(make([]byte, 20))
}

func TestNewBlockchain(t *testing.T) {
	creator := randomAddress()
	bc, err := NewBlockchain(50, 8, creator, sign_ecdsa.EcdsaSigner{})
	if err != nil {
		t.Fatalf("NewBlockchain failed: %v", err)
	}
	if len(bc.Blocks) != 1 {
		t.Errorf("Genesis block not created")
	}
	if bc.CurrentRewards != 50 {
		t.Errorf("CurrentRewards = %d, want 50", bc.CurrentRewards)
	}
	if bc.CurrentDifficult != 8 {
		t.Errorf("CurrentDifficult = %d, want 8", bc.CurrentDifficult)
	}
}

func TestAddTransactionToPool(t *testing.T) {
	creator := randomAddress()
	bc, _ := NewBlockchain(50, 8, creator, sign_ecdsa.EcdsaSigner{})
	signature := generateTestKeys(t, bc.signer)
	tx, _ := transaction.CreateTransaction(coin_transfer.CoinTransfer,creator, 10, 1, map[string]any{
		"recipient": randomAddress(),
	})
	tx.AddSing(bc.signer, signature)
	err := bc.AddTransactionToPool(tx)
	if err != nil {
		t.Fatalf("AddTransactionToPool failed: %v", err)
	}
	if len(bc.TxPool) != 1 {
		t.Errorf("TxPool should have 1 transaction")
	}
}

func TestMineBlockFromPool(t *testing.T) {
	creator := randomAddress()
	bc, _ := NewBlockchain(50, 8, creator, sign_ecdsa.EcdsaSigner{})
	signature := generateTestKeys(t, bc.signer)
	tx, _ := transaction.CreateTransaction(coin_transfer.CoinTransfer,creator, 10, 1, map[string]any{
		"recipient": randomAddress(),
	})
	tx.AddSing(bc.signer, signature)
	bc.AddTransactionToPool(tx)
	_, err := bc.MineBlockFromPool(creator)
	if err != nil {
		t.Fatalf("MineBlockFromPool failed: %v", err)
	}
	if len(bc.Blocks) != 2 {
		t.Errorf("Expected 2 blocks after mining, got %d", len(bc.Blocks))
	}
	if len(bc.TxPool) != 0 {
		t.Errorf("TxPool should be empty after mining")
	}
}

func TestAddBlock_InvalidBlock(t *testing.T) {
	creator := randomAddress()
	bc, _ := NewBlockchain(50, 8, creator, sign_ecdsa.EcdsaSigner{})
	b, _ := block.NewBlock(&bc.Blocks[len(bc.Blocks)-1], bc.CurrentDifficult)
	b.Hash[0] = 1 // Corrupt hash
	err := bc.AddBlock(b)
	if err == nil {
		t.Error("Expected error for invalid block hash")
	}
}

func TestVerifyBlockchain(t *testing.T) {
	creator := randomAddress()
	bc, _ := NewBlockchain(50, 8, creator, sign_ecdsa.EcdsaSigner{})
	signature := generateTestKeys(t, bc.signer)
	for i := 0; i < 3; i++ {
		tx, _ := transaction.CreateTransaction(coin_transfer.CoinTransfer,creator, int64(i+1), 1, map[string]any{
			"recipient": randomAddress(),
		})
		tx.AddSing(bc.signer, signature)
		bc.AddTransactionToPool(tx)
		bc.MineBlockFromPool(creator)
	}
	err := bc.Verify(4)
	if err != nil {
		t.Errorf("Blockchain verification failed: %v", err)
	}
}
