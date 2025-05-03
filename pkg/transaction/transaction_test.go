package transaction

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
	tx, err := NewTransaction(sender, reciver, 100)
	if err != nil {
		t.Fatalf("NewTransaction failed: %v", err)
	}
	if tx.Value != 100 {
		t.Errorf("Value = %d, want 100", tx.Value)
	}
	if hex.EncodeToString(tx.Sender[:]) != sender {
		t.Errorf("Sender mismatch")
	}
	if hex.EncodeToString(tx.Reciver[:]) != reciver {
		t.Errorf("Reciver mismatch")
	}
}

func TestNewTransaction_InvalidAddress(t *testing.T) {
	_, err := NewTransaction("abc", randomAddress(), 100)
	if err == nil {
		t.Error("Expected error for invalid sender address")
	}
	_, err = NewTransaction(randomAddress(), "abc", 100)
	if err == nil {
		t.Error("Expected error for invalid reciver address")
	}
}

func TestCalcTxId(t *testing.T) {
	tx, _ := NewTransaction(randomAddress(), randomAddress(), 42)
	tx.CalcTxId()
	if tx.TxId == [32]byte{} {
		t.Error("TxId should not be zero after CalcTxId")
	}
}

func TestSignAndVerify(t *testing.T) {
	signature := generateTestKeys(t)
	sender := randomAddress()
	reciver := randomAddress()
	tx, _ := NewTransaction(sender, reciver, 55)
	tx.CalcTxId()
	err := tx.Sing(signature)
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
	tx, _ := NewTransaction(randomAddress(), randomAddress(), 1)
	tx.CalcTxId()
	tx.Sing(signature)
	tx.Sign[0] ^= 0xFF // Corrupt signature
	err := tx.Verify()
	if err == nil {
		t.Error("Expected error for invalid signature")
	}
}

func TestVerify_TamperedData(t *testing.T) {
	signature := generateTestKeys(t)
	tx, _ := NewTransaction(randomAddress(), randomAddress(), 1)
	tx.CalcTxId()
	tx.Sing(signature)
	tx.Value = 999 // Tamper with transaction
	err := tx.Verify()
	if err == nil {
		t.Error("Expected error for tampered transaction data")
	}
}
