package ballance_storage

import (
	"testing"
)

func TestAddAndSubBallance(t *testing.T) {
	storage := NewMemoryStorage()
	address := "test_address"

	// Add balance
	newBal, err := storage.AddBallance(address, 100)
	if err != nil {
		t.Fatalf("AddBallance failed: %v", err)
	}
	if newBal != 100 {
		t.Errorf("Expected balance 100, got %d", newBal)
	}

	// Subtract balance
	newBal, err = storage.SubBallance(address, 40)
	if err != nil {
		t.Fatalf("SubBallance failed: %v", err)
	}
	if newBal != 60 {
		t.Errorf("Expected balance 60, got %d", newBal)
	}
}

func TestTransfer(t *testing.T) {
	storage := NewMemoryStorage()
	sender := "sender"
	receiver := "receiver"
	storage.AddBallance(sender, 200)
	storage.AddBallance(receiver, 50)

	err := storage.Transfer(sender, receiver, 70)
	if err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}

	if storage.GetBallance(sender) != 130 {
		t.Errorf("Expected sender balance 130, got %d", storage.GetBallance(sender))
	}
	if storage.GetBallance(receiver) != 120 {
		t.Errorf("Expected receiver balance 120, got %d", storage.GetBallance(receiver))
	}
}

func TestConfirmAndReject(t *testing.T) {
	storage := NewMemoryStorage()
	address := "test_address"
	storage.AddBallance(address, 100)

	// Confirm should move txPool to ballancePool
	err := storage.Confirm()
	if err != nil {
		t.Fatalf("Confirm failed: %v", err)
	}
	if storage.GetBallance(address) != 100 {
		t.Errorf("Expected balance 100 after confirm, got %d", storage.GetBallance(address))
	}

	// Add more, then reject
	storage.AddBallance(address, 50)
	err = storage.Reject()
	if err != nil {
		t.Fatalf("Reject failed: %v", err)
	}
	if storage.GetBallance(address) != 100 {
		t.Errorf("Expected balance 100 after reject, got %d", storage.GetBallance(address))
	}
}
