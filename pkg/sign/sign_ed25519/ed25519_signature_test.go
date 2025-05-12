package sign_ed25519

import (
	"crypto/sha256"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	signer := Ed25519Signer{}
	signature, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	if len(signature.PrivateKey) != 64 {
		t.Errorf("Private key length = %d, want 64", len(signature.PrivateKey))
	}
	if len(signature.PublicKey) != 32 {
		t.Errorf("Public key length = %d, want 32", len(signature.PublicKey))
	}

}

func TestSignAndVerify(t *testing.T) {
	signer := Ed25519Signer{}
	signature, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	msg := []byte("test message")
	hash := sha256.Sum256(msg)
	sig, err := signer.Sign(hash[:], signature.PrivateKey)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	valid, err := signer.Verify(hash[:], sig, signature.PublicKey)
	if err != nil {
		t.Fatalf("Sign (verify) failed: %v", err)
	}
	if !valid {
		t.Error("Signature should be valid, got false")
	}
}

func TestVerifyInvalidSignature(t *testing.T) {
	signer := Ed25519Signer{}
	signature, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	msg := []byte("test message")
	hash := sha256.Sum256(msg)
	sig, err := signer.Sign(hash[:], signature.PrivateKey)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	// Corrupt the signature
	sig[0] ^= 0xFF
	valid, err := signer.Verify(hash[:], sig, signature.PublicKey)
	if err != nil {
		t.Fatalf("Sign (verify) failed: %v", err)
	}
	if valid {
		t.Error("Corrupted signature should be invalid, got true")
	}
}

func TestVerifyWithWrongKey(t *testing.T) {
	signer := Ed25519Signer{}
	signature1, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	signature2, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	msg := []byte("test message")
	hash := sha256.Sum256(msg)
	sig, err := signer.Sign(hash[:], signature1.PrivateKey)
	if err != nil {
		t.Fatalf("Sing failed: %v", err)
	}
	valid, err := signer.Verify(hash[:], sig, signature2.PublicKey)
	if err != nil {
		t.Fatalf("Sign (verify) failed: %v", err)
	}
	if valid {
		t.Error("Signature verified with wrong public key, should be invalid")
	}
}
