package merkle

import (
	"blockchain_demo/pkg/transaction"
	"crypto/sha256"
	"fmt"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type MerkeTree struct {
	Tree   []transaction.Hash
	Levels []int64
}

func (merkle *MerkeTree) Count() int64 {
	if len(merkle.Levels) == 0 {
		return 0
	}

	return merkle.Levels[0]
}

func (merkle *MerkeTree) Root() transaction.Hash {
	if len(merkle.Tree) == 0 {
		return transaction.Hash{}
	}

	return merkle.Tree[len(merkle.Tree)-1]
}

func CreateMerkeTree(hashes []transaction.Hash) (*MerkeTree, error) {
	tree := make([]transaction.Hash, len(hashes))
	copy(tree, hashes)
	var merkle = MerkeTree{
		Tree:   tree,
		Levels: []int64{int64(len(hashes))},
	}
	merkle.rebuild()
	return &merkle, nil
}

func getHash(item1 transaction.Hash, item2 transaction.Hash) transaction.Hash {
	var hasher1 = sha256.New()
	hasher1.Write(item1[:])
	hasher1.Write(item2[:])

	var hasher2 = sha256.New()
	hasher2.Write(hasher1.Sum(nil))
	return transaction.Hash(hasher2.Sum(nil))
}

func (merkle *MerkeTree) buildNextLevel(current []transaction.Hash) []transaction.Hash {
	var level = []transaction.Hash{}
	for i := 0; i < len(current); i += 2 {
		var j = i + 1
		if i == len(current)-1 {
			j = i
		}
		level = append(level, getHash(current[i], current[j]))
	}
	return level
}
func (merkle *MerkeTree) rebuild() {
	var current = merkle.Tree[:merkle.Levels[0]]
	for len(current) > 1 {
		current = merkle.buildNextLevel(current)
		merkle.Tree = append(merkle.Tree, current...)
		if len(current) > 1 {
			merkle.Levels = append(merkle.Levels, int64(len(current)))
		}
	}
}

func (merkle *MerkeTree) GetMerkleProof(index int64) ([][33]byte, error) {
	if index < 0 || index >= merkle.Count() {
		return nil, fmt.Errorf("transaction with index %T not found", index)
	}

	var currentIndex = index
	var siblins = [][33]byte{}
	var offset int64 = 0
	for _, levelSize := range merkle.Levels {
		var leftRight = (currentIndex & 1)
		// -1 or +1 due to even or odd index to get left or right neighbour
		var shift = currentIndex + (1 - ((leftRight & 1) << 1))
		var levelIndex = offset + min(levelSize-1, shift)
		siblin := append([]byte{byte(leftRight)}, merkle.Tree[levelIndex][:]...)
		siblins = append(siblins, [33]byte(siblin))
		currentIndex >>= 1
		offset += levelSize
	}

	return siblins, nil
}

func VerifyMerkleProof(txHash transaction.Hash, siblings [][33]byte, root transaction.Hash) bool {
	var currentHash = txHash
	for _, sibling := range siblings {
		var leftRight = sibling[:1][0]
		var hash = sibling[1:][:]
		if leftRight == 0 {
			currentHash = getHash(currentHash, transaction.Hash(hash))
		} else {
			currentHash = getHash(transaction.Hash(hash), currentHash)
		}
	}
	return [32]byte(currentHash) == root
}
