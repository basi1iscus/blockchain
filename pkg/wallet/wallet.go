package wallet

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/utils"
	"bytes"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	Keys    sign.SignatureKeys
	Address string
}

func createAddress(pubKey []byte, prefix []byte) (string, error) {
	hashed, err := utils.GetHash(pubKey)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: %v", err)
	}

	hasher := ripemd160.New()
	hasher.Write(hashed)
	netAddress := hasher.Sum(prefix)

	hash, err := utils.GetHash(netAddress)
	if err != nil {
		return "", fmt.Errorf("failed to hash netAddress: %v", err)
	}
	checksum, err := utils.GetHash(hash)
	if err != nil {
		return "", fmt.Errorf("failed to hash checksum: %v", err)
	}
	address := append(netAddress, checksum[:4]...)

	return base58.Encode(address), nil
}

func CreateWallet(keys *sign.SignatureKeys, prefix []byte) (*Wallet, error) {
	address, err := createAddress(keys.PublicKey, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}

	return &Wallet{Keys: *keys, Address: address}, nil
}

func ValidateAddress(pubKey []byte, prefix []byte, address string) error {
	addr, err := createAddress(pubKey, prefix)
	if err != nil {
		return fmt.Errorf("failed to create address: %v", err)
	}

	if addr != address {
		return fmt.Errorf("address not match to public key")
	}

	return nil
}

func (w Wallet) GetPublicKeyHash() ([]byte, error) {
	hashed, err := utils.GetHash(w.Keys.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create hash: %v", err)
	}

	hasher := ripemd160.New()
	hasher.Write(hashed)
	netAddress := hasher.Sum(nil)
	return netAddress, nil
}

func CheckAddress(address string) error {
	binAddress := base58.Decode(address)
	
	hash, err := utils.GetHash(binAddress[:21])
	if err != nil {
		return fmt.Errorf("failed to hash netAddress: %v", err)
	}
	checksum, err := utils.GetHash(hash)
	if err != nil {
		return fmt.Errorf("failed to hash checksum: %v", err)
	}
	if !bytes.Equal(checksum[:4], binAddress[21:]) {
		return fmt.Errorf("address nor correct")
	}
	
	return nil
}
