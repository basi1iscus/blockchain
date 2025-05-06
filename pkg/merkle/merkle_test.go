package merkle

import (
	"blockchain_demo/pkg/transaction"
	"encoding/hex"
	"testing"
)

func mustHashFromHex(s string) transaction.Hash {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	var h transaction.Hash
	copy(h[:], b)
	return h
}

func TestMerkeTree_CreateMerkeTree(t *testing.T) {
	h1 := mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000001")
	h2 := mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000002")
	h3 := mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000003")

	tree, err := CreateMerkeTree([]transaction.Hash{h1, h2, h3})
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	if tree.Count() != 3 {
		t.Errorf("Expected count 3, got %d", tree.Count())
	}

}

func TestMerkeTree_Count(t *testing.T) {
	merkle, _ := CreateMerkeTree([]transaction.Hash{
		mustHashFromHex("01"),
		mustHashFromHex("02"),
	})
	if merkle.Count() != 2 {
		t.Errorf("Expected count 2, got %d", merkle.Count())
	}
}

func TestMerkeTree_Root(t *testing.T) {
	merkle, _ := CreateMerkeTree([]transaction.Hash{
		mustHashFromHex("01"),
		mustHashFromHex("02"),
	})
	root := merkle.Root()
	if root == (transaction.Hash{}) {
		t.Errorf("Expected non-empty root")
	}
}

func TestMerkeTree_GetMerkleProof_and_VerifyMerkleProof(t *testing.T) {
	hashes := []transaction.Hash{
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000001"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000002"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000003"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000004"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000005"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000006"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000007"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000008"),
		mustHashFromHex("0000000000000000000000000000000000000000000000000000000000000009"),
		mustHashFromHex("000000000000000000000000000000000000000000000000000000000000000A"),
		mustHashFromHex("000000000000000000000000000000000000000000000000000000000000000B"),
		mustHashFromHex("000000000000000000000000000000000000000000000000000000000000000C"),
	}
	merkle, _ := CreateMerkeTree(hashes)
	for i, h := range hashes {
		proof, err := merkle.GetMerkleProof(int64(i))
		if err != nil {
			t.Fatalf("GetMerkleProof failed: %v", err)
		}
		if !VerifyMerkleProof(h, proof, merkle.Root()) {
			t.Errorf("Merkle proof verification failed for index %d", i)
		}
	}
}
