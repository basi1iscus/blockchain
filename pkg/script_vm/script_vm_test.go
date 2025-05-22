package script_vm

import (
	"blockchain_demo/pkg/sign/sign_ed25519"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"blockchain_demo/pkg/wallet"
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func TestVM_P2PKH(t *testing.T) {
	// 1. Generate key pair
	signer := sign_ed25519.Ed25519Signer{}
	keys, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	// 2. Prepare transaction (dummy, just for txid)
	prefix := []byte{0x00} // Example prefix for mainnet
	wallet, err := wallet.CreateWallet(keys, prefix)
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}
	pubKeyHash, err := wallet.GetPublicKeyHash()
	address := hex.EncodeToString(pubKeyHash)
	if err != nil {
		t.Fatalf("failed to get public key hash: %v", err)
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx, err := transaction.CreateTransaction(coin_transfer.CoinTransfer, address, rnd.Int63n(1000), rnd.Int63n(10), map[string]any{
		"recipient": address,
	})
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}
	// 3. Prepare signature (sign txid)
	txid := tx.GetTxId()
	sig, err := signer.Sign(txid[:], keys.PrivateKey)
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	// 4. ScriptSig: <sig> <pubkey>
	scriptSig := append(
		append([]byte{byte(len(sig))}, sig...),
		append([]byte{byte(len(keys.PublicKey))}, keys.PublicKey...)...,
	)

	// 5. ScriptPubKey: OP_DUP OP_HASH160 <pubkeyhash> OP_EQUALVERIFY OP_CHECKSIG

	if err != nil {
		t.Fatalf("failed to get public key hash: %v", err)
	}
	scriptPubKey := []byte{
		OP_DUP,
		OP_HASH160,
		byte(len(pubKeyHash)),
	}
	scriptPubKey = append(scriptPubKey, pubKeyHash...)
	scriptPubKey = append(scriptPubKey, OP_EQUALVERIFY, OP_CHECKSIG)

	// 6. Concatenate scripts (as in Bitcoin: ScriptSig || ScriptPubKey)
	fullScript := append(scriptSig, scriptPubKey...)

	// 7. Run VM
	vm := New(&signer)
	err = vm.Execute(fullScript, tx)
	if err != nil {
		t.Fatalf("VM failed: %v", err)
	}

	// 8. Check result (should be OP_TRUE on stack)
	result, err := vm.stack.Pop()
	if err != nil {
		t.Fatalf("stack empty: %v", err)
	}
	if !bytes.Equal(result, []byte{OP_TRUE}) {
		t.Errorf("expected OP_TRUE, got %x", result)
	}
}
