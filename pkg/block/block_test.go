package block

import (
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"testing"
)

func TestNewBlock_NoPrev(t *testing.T) {
	block, err := NewBlock(nil, 8)
	if err != nil {
		t.Fatalf("NewBlock failed: %v", err)
	}
	if block.Index != 1 {
		t.Errorf("Index = %d, want 1", block.Index)
	}
	if block.Difficulty != 8 {
		t.Errorf("Difficulty = %d, want 8", block.Difficulty)
	}
	if block.Prev != [32]byte{} {
		t.Errorf("Prev should be zeroed")
	}
}

func TestNewBlock_WithPrev(t *testing.T) {
	prev, _ := NewBlock(nil, 8)
	prev.Index = 5
	prev.Hash = [32]byte{1, 2, 3}
	block, err := NewBlock(prev, 16)
	if err != nil {
		t.Fatalf("NewBlock failed: %v", err)
	}
	if block.Index != 6 {
		t.Errorf("Index = %d, want 6", block.Index)
	}
	if block.Prev != prev.Hash {
		t.Errorf("Prev not set to previous block's hash")
	}
}

func TestCalcHash(t *testing.T) {
	block, _ := NewBlock(nil, 8)
	hash, _ := block.CalcHash(uint64(0))
	if len(hash) != 32 {
		t.Errorf("Hash length = %d, want 32", len(hash))
	}
}

func TestMine(t *testing.T) {
	block, _ := NewBlock(nil, 16) // Low difficulty for test
	hash, _ := block.Mine(0)
	if len(hash) != 32 {
		t.Errorf("Hash length = %d, want 32", len(hash))
	}
	if block.Hash != [32]byte(hash) {
		t.Errorf("Block.Hash not set correctly after mining")
	}
	if block.Nonce == 0 {
		t.Errorf("Nonce should be non-zero after mining")
	}
}

func TestMineHalfByteDifficult(t *testing.T) {
	block, _ := NewBlock(nil, 12) // Low difficulty for test
	hash, _ := block.Mine(0)
	if len(hash) != 32 {
		t.Errorf("Hash length = %d, want 32", len(hash))
	}
	if block.Hash != [32]byte(hash) {
		t.Errorf("Block.Hash not set correctly after mining")
	}
	if block.Nonce == 0 {
		t.Errorf("Nonce should be non-zero after mining")
	}
}

func TestBlockWithTransactions(t *testing.T) {
	block, _ := NewBlock(nil, 8)
	tx, _ := transaction.CreateTransaction(coin_transfer.CoinTransfer, "00112233445566778899aabbccddeeff00112233", 10, 1, map[string]any{
		"recipient": "ffeeddccbbaa99887766554433221100ffeeddcc",
	})
	block.Transactions = append(block.Transactions, tx)
	hash, _ := block.Mine(0)
	if len(hash) != 32 {
		t.Errorf("Hash length = %d, want 32", len(hash))
	}
}
