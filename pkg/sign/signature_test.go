package sign

import (
	"crypto/sha256"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	signature, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	if len(signature.PrivateKey) != 121 {
		t.Errorf("Private key length = %d, want 64", len(signature.PrivateKey))
	}
	if len(signature.PublicKey) != 91 {
		t.Errorf("Public key length = %d, want 91", len(signature.PublicKey))
	}

}

func TestSignAndVerify(t *testing.T) {
	signature, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	msg := []byte("test message")
	hash := sha256.Sum256(msg)
	sig, err := Sign(hash[:], signature.PrivateKey)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	valid, err := Verify(hash[:], sig, signature.PublicKey)
	if err != nil {
		t.Fatalf("Sign (verify) failed: %v", err)
	}
	if !valid {
		t.Error("Signature should be valid, got false")
	}
}

func TestVerifyInvalidSignature(t *testing.T) {
	signature, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	msg := []byte("test message")
	hash := sha256.Sum256(msg)
	sig, err := Sign(hash[:], signature.PrivateKey)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	// Corrupt the signature
	sig[0] ^= 0xFF
	valid, err := Verify(hash[:], sig, signature.PublicKey)
	if err != nil {
		t.Fatalf("Sign (verify) failed: %v", err)
	}
	if valid {
		t.Error("Corrupted signature should be invalid, got true")
	}
}

func TestVerifyWithWrongKey(t *testing.T) {
	signature1, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	signature2, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	msg := []byte("test message")
	hash := sha256.Sum256(msg)
	sig, err := Sign(hash[:], signature1.PrivateKey)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	valid, err := Verify(hash[:], sig, signature2.PublicKey)
	if err != nil {
		t.Fatalf("Sign (verify) failed: %v", err)
	}
	if valid {
		t.Error("Signature verified with wrong public key, should be invalid")
	}
}
