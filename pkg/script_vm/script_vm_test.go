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
	err = vm.Run(fullScript, tx)
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

func TestVM_AllOpcodes(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	keys, err := signer.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	prefix := []byte{0x00}
	w, err := wallet.CreateWallet(keys, prefix)
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}
	pubKeyHash, _ := w.GetPublicKeyHash()
	tx, _ := transaction.CreateTransaction(coin_transfer.CoinTransfer, hex.EncodeToString(pubKeyHash), 1, 1, map[string]any{"recipient": hex.EncodeToString(pubKeyHash)})
	txid := tx.GetTxId()
	sig, _ := signer.Sign(txid[:], keys.PrivateKey)
	pub := keys.PublicKey

	cases := []struct {
		name   string
		script []byte
		want   []byte
		wantErr bool
		}{
		{"OP_0", []byte{OP_0}, []byte{OP_0}, false},
		{"OP_1NEGATE", []byte{OP_1NEGATE}, []byte{OP_1NEGATE}, false},
		{"OP_1", []byte{OP_1}, []byte{1}, false},
		{"OP_16", []byte{OP_16}, []byte{16}, false},
		{"OP_PUSHDATA0_01", append([]byte{1}, []byte{0xAB}...), []byte{0xAB}, false},
		{"OP_PUSHDATA1", []byte{OP_PUSHDATA1, 1, 0xCD}, []byte{0xCD}, false},
		{"OP_DUP", []byte{1, 0xAA, OP_DUP}, []byte{0xAA}, false},
		{"OP_DROP", []byte{1, 0xAA, OP_DROP}, nil, false},
		{"OP_IFDUP_nonzero", []byte{1, 0x01, OP_IFDUP}, []byte{0x01}, false},
		{"OP_IFDUP_zero", []byte{1, 0x00, OP_IFDUP}, []byte{0x00}, false},
		{"OP_EQUAL_true", []byte{1, 0x01, 1, 0x01, OP_EQUAL}, []byte{OP_TRUE}, false},
		{"OP_EQUAL_false", []byte{1, 0x01, 1, 0x02, OP_EQUAL}, []byte{OP_FALSE}, false},
		{"OP_EQUALVERIFY_true", []byte{1, 0x01, 1, 0x01, OP_EQUALVERIFY}, nil, false},
		{"OP_EQUALVERIFY_false", []byte{1, 0x01, 1, 0x02, OP_EQUALVERIFY}, nil, true},
		{"OP_VERIFY_true", []byte{1, 0x01, OP_VERIFY}, nil, false},
		{"OP_VERIFY_false", []byte{1, 0x00, OP_VERIFY}, nil, true},
		{"OP_SHA256", []byte{1, 0x01, OP_SHA256}, nil, false},
		{"OP_HASH160", []byte{1, 0x01, OP_HASH160}, nil, false},
		{"OP_HASH256", []byte{1, 0x01, OP_HASH256}, nil, false},
		{"OP_CHECKSIG_true", append(append(append([]byte{byte(len(sig))}, sig...), append([]byte{byte(len(pub))}, pub...)...), OP_CHECKSIG), []byte{OP_TRUE}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vm := New(&signer)
			err := vm.Run(tc.script, tx)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.want != nil {
				res, err := vm.stack.Pop()
				if err != nil {
					t.Errorf("stack error: %v", err)
				}
				if !bytes.Equal(res, tc.want) {
					t.Errorf("expected %x, got %x", tc.want, res)
				}
			}
		})
	}
}
