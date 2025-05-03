package sign

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"
)

type Signature struct {
	PrivateKey []byte
	PublicKey  []byte
}

// GenerateKeyPair generates a new ECDSA private and public key pair for elliptic.P256().
// Returns the private key and the public key as hex string.
func GenerateKeyPair() (*Signature, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	x509Encoded, _ := x509.MarshalECPrivateKey(privKey)
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(&privKey.PublicKey)

	return &Signature{PrivateKey: x509Encoded, PublicKey: x509EncodedPub}, nil
}

// privateKey is x509Encoded hex string
func Sign(data []byte, privateKey []byte) ([]byte, error) {
	// Decode the private key from hex string
	privKey, _ := x509.ParseECPrivateKey(privateKey)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, data[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign")
	}

	// Serialize the signature (r, s) into a 65-byte array
	signature := make([]byte, 65)
	copy(signature[:32], r.Bytes())
	copy(signature[32:64], s.Bytes())
	signature[64] = 0 // Recovery ID placeholder (not used here)

	return signature, nil
}

// publicKey is x509 PKIX Encoded hex string
func Verify(data []byte, signature []byte, publicKey []byte) (bool, error) {
	pubKey, _ := x509.ParsePKIXPublicKey(publicKey)
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("invalid public key type")
	}

	// Parse the signature into r and s
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Verify the signature
	isValid := ecdsa.Verify(ecdsaPubKey, data, r, s)
	return isValid, nil
}
