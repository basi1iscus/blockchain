package transaction_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"blockchain_demo/pkg/transaction"
	_ "blockchain_demo/pkg/transaction/coin_transfer"
	_ "blockchain_demo/pkg/transaction/contract_call"
	_ "blockchain_demo/pkg/transaction/token_transfer"
)

func TestTransaction_SerializeDeserialize(t *testing.T) {

	sender := "1234567890abcdef1234567890abcdef12345678"
	address := "1234567890abcdef1234567890abcdef12345678"
	
	types := []struct {
		typeName transaction.TransactionType
		params  map[string]any
	}{
		{
			typeName: "coin_transfer",
			params: map[string]any{
				"recipient": address,
			},
		},
		{
			typeName: "contract_call",
			params: map[string]any{
				"contractAddress": address,
				"to": address,
				"contractType": "transfer",
				"amount": uint64(42),
			},
		},
		{
			typeName: "token_transfer",
			params: map[string]any{
				"token": address,
				"recipient": address,
				"amount": int64(123),
			},
		},
	}

	value := int64(100)
	fee := int64(1)

	for _, typ := range types {
		tx, err := transaction.CreateTransaction(typ.typeName, sender, value, fee, typ.params)
		if err != nil {
			t.Fatalf("failed to create %s transaction: %v", typ.typeName, err)
		}

		data, err := tx.Stringify()
		if err != nil {
			t.Fatalf("failed to serialize %s transaction: %v", typ.typeName, err)
		}

		deserialized, err := transaction.ParseTransaction(data)
		if err != nil {
			t.Fatalf("failed to deserialize %s transaction: %v", typ.typeName, err)
		}

		// Compare types
		if reflect.TypeOf(tx) != reflect.TypeOf(deserialized) {
			t.Errorf("type mismatch after deserialization: got %T, want %T", deserialized, tx)
		}

		// Compare JSON representations
		origJSON, _ := json.Marshal(tx)
		desJSON, _ := json.Marshal(deserialized)
		if string(origJSON) != string(desJSON) {
			t.Errorf("json mismatch after deserialization for %s: got %s, want %s", typ.typeName, desJSON, origJSON)
		}
	}
}
