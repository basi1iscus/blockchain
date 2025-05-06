package block

import (
	"blockchain_demo/pkg/merkle"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/utils"
	"context"
	"fmt"
	"runtime"
	"time"
)

type Block struct {
	Index        uint32
	Time         int64
	Hash         [32]byte
	Prev         [32]byte
	Nonce        uint64
	Difficulty   uint64
	MerkleRoot   [32]byte
	Transactions []transaction.Transaction
}

func NewBlock(prevBlock *Block, difficulty uint64) (*Block, error) {
	var index uint32 = 1
	var prev [32]byte
	if prevBlock != nil {
		index = prevBlock.Index + 1
		prev = prevBlock.Hash
	}
	var block = Block{
		Index:        index,
		Time:         time.Now().UnixNano(),
		Hash:         [32]byte{},
		Prev:         prev,
		Nonce:        0,
		Difficulty:   difficulty,
		MerkleRoot:   [32]byte{},
		Transactions: []transaction.Transaction{},
	}

	return &block, nil
}

func (block *Block) CalcHash(nonce uint64) ([]byte, error) {
	var hash, err = utils.GetHash(block.Index, block.Time, block.Prev[:], nonce, block.MerkleRoot[:])
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func miner(block *Block, from uint64, count uint64, ch chan uint64, ctx context.Context) {

	var bytes uint64 = block.Difficulty / 8
	var bits uint64 = block.Difficulty % 8
	var buf = make([]byte, bytes)
	for n := uint64(0); n < bytes; n++ {
		buf[n] = 255
	}
	if bits > 0 {
		buf = append(buf, (255 << bits))
	}

nonceSearch:
	for nonce := from; nonce < from+count && ctx.Err() == nil; nonce++ {
		var hash, _ = block.CalcHash(nonce)
		for n, v := range buf {
			if hash[n]&v != 0 {
				continue nonceSearch
			}
		}
		ch <- nonce
		break
	}
}

func (block *Block) calcMerkleRoot() (*transaction.Hash, error) {
	var txHashes []transaction.Hash
	for _, tx := range block.Transactions {
		var txHash = tx.GetTxId()
		txHashes = append(txHashes, txHash)
	}
	var tree, err = merkle.CreateMerkeTree(txHashes)
	if err != nil {
		return nil, fmt.Errorf("error creating Merkle Tree: %s", err)
	}
	root := tree.Root()
	return &root, nil
}

func (block *Block) Mine(threads uint64) ([]byte, error) {
	root, err := block.calcMerkleRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate Merkle root: %w", err)
	}
	block.MerkleRoot = *root
	channel := make(chan uint64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if threads == 0 {
		threads = uint64(runtime.NumCPU())
	}
	for th := uint64(0); th < threads; th++ {
		var count = ^uint64(0) / threads
		go miner(block, count*th, count, channel, ctx)
	}
	var nonce = <-channel
	var hash, errHash = block.CalcHash(nonce)
	if errHash != nil {
		return nil, errHash
	}
	block.Nonce = nonce
	block.Hash = [32]byte(hash)
	return hash, nil
}

func (block *Block) Verify() error {
	var hash, err = block.CalcHash(block.Nonce)
	if err != nil {
		return err
	}
	if block.Hash != [32]byte(hash) {
		return fmt.Errorf("Block hash is invalid")
	}
	for _, tx := range block.Transactions {
		var err = tx.Verify()
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
	}
	root, err := block.calcMerkleRoot()
	if err != nil {
		return err
	}
	if *root != block.MerkleRoot {
		return fmt.Errorf("merkle root is invalid: %x", *root)
	}
	return nil
}

func (block *Block) AddTransaction(tx *transaction.Transaction) error {
	block.Transactions = append(block.Transactions, *tx)

	return nil
}

func (b *Block) String() string {
	return fmt.Sprintf("Block{Index: %d, Time: %d, Hash: %x, Merkle root %x,  Prev: %x, Nonce: %d, Difficulty: %d, TxCount: %d}",
		b.Index, b.Time, b.Hash, b.MerkleRoot, b.Prev, b.Nonce, b.Difficulty, len(b.Transactions))
}
