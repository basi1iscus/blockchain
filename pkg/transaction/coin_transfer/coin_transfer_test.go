package coin_transfer

import (
	"blockchain_demo/pkg/sign"
	"encoding/hex"
	"testing"
)

func generateTestKeys(t *testing.T) *sign.Signature {
	signature, err := sign.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	return signature
}

func randomAddress() string {
	return hex.EncodeToString(make([]byte, 20))
}

func TestNewTransaction_Valid(t *testing.T) {
	sender := randomAddress()
	reciver := randomAddress()
	tx, err := NewTransaction(sender, 100, 1, map[string]any{
		"recipient": reciver,
	})
	if err != nil {
		t.Fatalf("NewTransaction failed: %v", err)
	}
	if tx.Value != 100 {
		t.Errorf("Value = %d, want 100", tx.Value)
	}
	if hex.EncodeToString(tx.Sender[:]) != sender {
		t.Errorf("Sender mismatch")
	}
	if hex.EncodeToString(tx.Recipient[:]) != reciver {
		t.Errorf("Reciver mismatch")
	}
}

func TestNewTransaction_InvalidAddress(t *testing.T) {
	_, err := NewTransaction("abc", 100, 1, map[string]any{
		"recipient": randomAddress(),
	})
	if err == nil {
		t.Error("Expected error for invalid sender address")
	}
	_, err = NewTransaction(randomAddress(), 100, 1, map[string]any{
		"recipient": 100,
	})
	if err == nil {
		t.Error("Expected error for invalid reciver address")
	}
}

func TestCalcTxId(t *testing.T) {
	tx, _ := NewTransaction(randomAddress(), 42, 1, map[string]any{
		"recipient": randomAddress(),
	})
	if tx.TxId == [32]byte{} {
		t.Error("TxId should not be zero after CalcTxId")
	}
}

func TestSignAndVerify(t *testing.T) {
	signature := generateTestKeys(t)
	sender := randomAddress()
	reciver := randomAddress()
	tx, _ := NewTransaction(sender, 55, 1, map[string]any{
		"recipient": reciver,
	})
	err := tx.AddSing(signature)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	err = tx.Verify()
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
}

func TestVerify_InvalidSignature(t *testing.T) {
	signature := generateTestKeys(t)
	tx, _ := NewTransaction(randomAddress(), 1, 0, map[string]any{
		"recipient": randomAddress(),
	})
	tx.AddSing(signature)
	tx.Sign[0] ^= 0xFF // Corrupt signature
	err := tx.Verify()
	if err == nil {
		t.Error("Expected error for invalid signature")
	}
}

func TestVerify_TamperedData(t *testing.T) {
	signature := generateTestKeys(t)
	tx, _ := NewTransaction(randomAddress(), 1, 0, map[string]any{
		"recipient": randomAddress(),
	})
	tx.AddSing(signature)
	tx.Value = 999 // Tamper with transaction
	err := tx.Verify()
	if err == nil {
		t.Error("Expected error for tampered transaction data")
	}
}
