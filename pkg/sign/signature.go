package sign

type SignatureKeys struct {
	PrivateKey []byte
	PublicKey  []byte
}

type Signer interface {
	GenerateKeyPair() (*SignatureKeys, error)
	Sign(data []byte, privateKey []byte) ([]byte, error)
    Verify(data []byte, signature []byte, publicKey []byte) (bool, error)
}
