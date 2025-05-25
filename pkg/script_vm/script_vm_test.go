package script_vm

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/sign/sign_ed25519"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"blockchain_demo/pkg/utils"
	"blockchain_demo/pkg/wallet"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func PrepareVMForP2PKH(t *testing.T, signer sign.Signer) ([]byte, transaction.Transaction) {
	// 1. Generate key pair
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

	return fullScript, tx
}

func TestVM_P2PKH(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	fullScript, tx := PrepareVMForP2PKH(t, signer) 
	// 7. Run VM
	vm := New(&signer)
	err := vm.ParseScript(fullScript)
	if err != nil {
		t.Fatalf("failed to parse script: %v", err)
	}
	signedData := tx.GetTxId()
	fmt.Printf("Check hash: %#x\n", signedData)
	fmt.Println("Parsed script:")
	fmt.Println(vm)
	_, err = vm.Execute(signedData[:], nil)
	if err != nil {
		t.Fatalf("VM failed: %v", err)
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
			signedData := tx.GetTxId()
			res, err := vm.Run(tc.script, signedData[:])
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if tc.want != nil {
				if !bytes.Equal(res, tc.want) {
					t.Errorf("expected %x, got %x", tc.want, res)
				}
			}
		})
	}
}

func TestVM_IfElse(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	// keys, err := signer.GenerateKeyPair() // not needed for these tests
	tx, _ := utils.GetHash([]byte("dummy transaction for testing if/else"))

	cases := []struct {
		name   string
		script []byte
		want   []byte
		wantErr bool
	}{
		{
			"OP_IF true branch",
			[]byte{1, 0x01, OP_IF, 1, 0xAA, OP_ELSE, 1, 0xBB, OP_ENDIF},
			[]byte{0xAA}, false,
		},
		{
			"OP_IF false branch",
			[]byte{1, 0x00, OP_IF, 1, 0xAA, OP_ELSE, 1, 0xBB, OP_ENDIF},
			[]byte{0xBB}, false,
		},
		{
			"OP_NOTIF true branch",
			[]byte{1, 0x00, OP_NOTIF, 1, 0xCC, OP_ELSE, 1, 0xDD, OP_ENDIF},
			[]byte{0xCC}, false,
		},
		{
			"OP_NOTIF false branch",
			[]byte{1, 0x01, OP_NOTIF, 1, 0xCC, OP_ELSE, 1, 0xDD, OP_ENDIF},
			[]byte{0xDD}, false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vm := New(&signer)
			res, err := vm.Run(tc.script, tx)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if tc.want != nil {
				if !bytes.Equal(res, tc.want) {
					t.Errorf("expected %x, got %x", tc.want, res)
				}
			}
		})
	}
}

func TestVM_CheckMultiSig(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	// Generate 3 key pairs
	keys := make([]*sign.SignatureKeys, 3)
	for i := 0; i < 3; i++ {
		k, err := signer.GenerateKeyPair()
		if err != nil {
			t.Fatalf("failed to generate key pair: %v", err)
		}
		keys[i] = k
	}


	// Generate a random 20-byte hex address for recipient
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	addrBytes := make([]byte, 20)
	rnd.Read(addrBytes)
	address := hex.EncodeToString(addrBytes)

	tx, err := transaction.CreateTransaction(coin_transfer.CoinTransfer, address, 1000, 10, map[string]any{
		"recipient": address,
	})
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}
		txid := tx.GetTxId()

	// Create 2 valid signatures (for keys[0] and keys[2])
	sig0, _ := signer.Sign(txid[:], keys[0].PrivateKey)
	sig2, _ := signer.Sign(txid[:], keys[2].PrivateKey)

	// Script: <2> <sig0> <sig2> <3> <pub0> <pub1> <pub2> OP_CHECKMULTISIG
	script := []byte{
		OP_0, // need 2 signatures
		byte(len(sig0)),
	}
	script = append(script, sig0...)
	script = append(script, byte(len(sig2)))
	script = append(script, sig2...)
	script = append(script, OP_2) // need 2 signatures from 3 pubkeys
	for _, k := range keys {
		script = append(script, byte(len(k.PublicKey)))
		script = append(script, k.PublicKey...)
	}
	script = append(script, OP_3) // 3 pubkeys
	// script = append(script, 0) // dummy (bitcoin bug)
	script = append(script, OP_CHECKMULTISIG)

	vm := New(&signer)
	signedData := tx.GetTxId()
	_, err = vm.Run(script, signedData[:])
	if err != nil {
		t.Fatalf("VM failed: %v", err)
	}
}
func TestVM_ParseString_BasicOpcodes(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	vm := New(&signer)

	cases := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		{
			"OP_0",
			"OP_0",
			[]byte{OP_0},
			false,
		},
		{
			"OP_1NEGATE",
			"OP_1NEGATE",
			[]byte{OP_1NEGATE},
			false,
		},
		{
			"OP_1",
			"OP_1",
			[]byte{OP_1},
			false,
		},
		{
			"OP_16",
			"OP_16",
			[]byte{OP_16},
			false,
		},
		{
			"OP_PUSHDATA with hex data",
			fmt.Sprintf("OP_PUSHDATA %x", []byte{0xAB, 0xCD}),
			append([]byte{2}, []byte{0xAB, 0xCD}...),
			false,
		},
		{
			"OP_DUP",
			"OP_DUP",
			[]byte{OP_DUP},
			false,
		},
		{
			"OP_HASH160",
			"OP_HASH160",
			[]byte{OP_HASH160},
			false,
		},
		{
			"OP_EQUALVERIFY",
			"OP_EQUALVERIFY",
			[]byte{OP_EQUALVERIFY},
			false,
		},
		{
			"OP_CHECKSIG",
			"OP_CHECKSIG",
			[]byte{OP_CHECKSIG},
			false,
		},
		{
			"Unknown opcode",
			"OP_UNKNOWN",
			nil,
			true,
		},
		{
			"Invalid hex data",
			"OP_PUSHDATA ZZZZ",
			nil,
			true,
		},
		{
			"Comment and whitespace",
			`
			# This is a comment
			OP_1
			# Another comment
			OP_DUP
			`,
			[]byte{OP_1, OP_DUP},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := vm.ParseString(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !bytes.Equal(got, tc.expected) {
				t.Errorf("expected %x, got %x", tc.expected, got)
			}
		})
	}
}

func TestVM_ParseString_IfElseBranch(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	vm := New(&signer)

	scriptStr := `
		OP_1
		OP_IF
			OP_2
		OP_ELSE
			OP_3
		OP_ENDIF
	`
	scriptBytes, err := vm.ParseString(scriptStr)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	// Should parse to: OP_1 OP_IF OP_2 OP_ELSE OP_3 OP_ENDIF
	expected := []byte{OP_1, OP_IF, OP_2, OP_ELSE, OP_3, OP_ENDIF}
	if !bytes.Equal(scriptBytes, expected) {
		t.Errorf("expected %x, got %x", expected, scriptBytes)
	}
}

func TestVM_ParseString_CompilePushdata1(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	vm := New(&signer)
	// Data > 0x4B triggers OP_PUSHDATA1
	data := make([]byte, 0x4C)
	for i := range data {
		data[i] = byte(i)
	}
	scriptStr := fmt.Sprintf("OP_PUSHDATA %x", data)
	scriptBytes, err := vm.ParseString(scriptStr)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	// Should start with OP_PUSHDATA1, length, then data
	if len(scriptBytes) < 2+len(data) {
		t.Fatalf("script too short")
	}
	if scriptBytes[0] != OP_PUSHDATA1 || scriptBytes[1] != byte(len(data)) {
		t.Errorf("expected OP_PUSHDATA1 and length, got %x %x", scriptBytes[0], scriptBytes[1])
	}
	if !bytes.Equal(scriptBytes[2:], data) {
		t.Errorf("data mismatch")
	}
}

func TestVM_String_ParseString_RoundTrip(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}
	vm := New(&signer)
	// Compose a script with various opcodes and data
	pubKeyHash := make([]byte, 20)
	for i := range pubKeyHash {
		pubKeyHash[i] = byte(i + 1)
	}
	script := []byte{
		OP_DUP,
		OP_HASH160,
		byte(len(pubKeyHash)),
	}
	script = append(script, pubKeyHash...)
	script = append(script, OP_EQUALVERIFY, OP_CHECKSIG)

	// Parse the script into VM
	err := vm.ParseScript(script)
	if err != nil {
		t.Fatalf("ParseScript failed: %v", err)
	}
	// Get string representation
	str := vm.String()

	// Now parse back using ParseString
	vm2 := New(&signer)
	parsed, err := vm2.ParseString(str)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	if !bytes.Equal(parsed, script) {
		t.Errorf("expected %x, got %x", script, parsed)
	}
}

func TestVM_Parse_And_Compile(t *testing.T) {
	signer := sign_ed25519.Ed25519Signer{}	

	script, tx := PrepareVMForP2PKH(t, signer) 
	// 7. Run VM
	vm := New(&signer)
	err := vm.ParseScript(script)
	if err != nil {
		t.Fatalf("failed to parse script: %v", err)
	}
	// Get string representation
	str := vm.String()
	vm2 := New(&signer)
	parsed, err := vm2.ParseString(str)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	if !bytes.Equal(parsed, script) {
		t.Errorf("expected %x, got %x", script, parsed)
	}

	fmt.Println("Parsed script:")
	fmt.Println(vm2)
	signedData := tx.GetTxId()
	err = vm2.ParseScript(parsed)
	if err != nil {	
		t.Fatalf("failed to parse script: %v", err)
	}
	_, err = vm2.Execute(signedData[:], nil)
	if err != nil {
		t.Fatalf("VM failed: %v", err)
	}	
}
