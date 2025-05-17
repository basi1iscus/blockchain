package wallet

import (
	"blockchain_demo/pkg/sign/sign_ed25519"
	"testing"
)

func TestGreateWallet(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	keys, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	prefix := []byte{0x00} // Example prefix for mainnet
	wallet, err := CreateWallet(keys, prefix)
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}

	if wallet.Address == "" {
		t.Error("wallet address should not be empty")
	}
}

func TestValidateAddress(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	keys, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	prefix := []byte{0x00} // Example prefix for mainnet
	wallet, err := CreateWallet(keys, prefix)
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}

	err = ValidateAddress(keys.PublicKey, prefix, wallet.Address)
	if err != nil {
		t.Errorf("address validation failed: %v", err)
	}

	// Test with an invalid address
	invalidAddress := "invalidAddress"
	err = ValidateAddress(keys.PublicKey, prefix, invalidAddress)
	if err == nil {
		t.Error("expected validation to fail for an invalid address")
	}
}
