package sign_ed25519

import (
	"blockchain_demo/pkg/sign"
	"crypto/ed25519"
	"fmt"
)

type Ed25519Signer struct {
}

// GenerateKeyPair generates a new ECDSA private and public key pair for elliptic.P256().
// Returns the private key and the public key
func (signer Ed25519Signer) GenerateKeyPair() (*sign.SignatureKeys, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return &sign.SignatureKeys{PrivateKey: privKey, PublicKey: pubKey}, nil
}

func (signer Ed25519Signer) Sign(data []byte, privateKey []byte) ([]byte, error) {
	// Decode the private key from hex string
	privKey := ed25519.PrivateKey(privateKey)
	sign := ed25519.Sign(privKey, data[:])

	// Serialize the signature (r, s) into a 65-byte array
	signature := make([]byte, 65)
	copy(signature[:64], sign)
	signature[64] = 0 // Recovery ID placeholder (not used here)

	return signature, nil
}

func (signer Ed25519Signer) Verify(data []byte, signature []byte, publicKey []byte) (bool, error) {
	pubKey := ed25519.PublicKey(publicKey)

	// Parse the signature into r and s
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}

	// Verify the signature
	isValid := ed25519.Verify(pubKey, data, signature[:64])
	return isValid, nil
}
