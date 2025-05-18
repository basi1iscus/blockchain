package blockchain

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/block"
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"blockchain_demo/pkg/transaction_processor"
	"fmt"
	"strings"
	"sync"
)

const EmptyAddress = "0000000000000000000000000000000000000000"

type Blockchain struct {
	CurrentDifficulty uint64
	CurrentRewards    uint64
	txPool            []transaction.Transaction
	blocks            []block.Block
	signature         *sign.SignatureKeys
	signer            sign.Signer
	txProcessor       transaction_processor.TransactionProcessor
	storage           ballance_storage.BallanceStorage
	mu                sync.Mutex
}

func NewBlockchain(rewards uint64, difficulty uint64, creator string, signer sign.Signer, storage ballance_storage.BallanceStorage, txTypes map[transaction.TransactionType]transaction_processor.TransactionProcessor) (*Blockchain, error) {
	var blockchain = Blockchain{
		CurrentDifficulty: difficulty,
		CurrentRewards:    rewards,
		txPool:            []transaction.Transaction{},
		blocks:            []block.Block{},
		signer:            signer,
		storage:           storage,
		txProcessor:       transaction_processor.BaseProcessor{},
	}

	for txType, processor := range txTypes {
		transaction_processor.RegisterProcessorType(txType, processor)
	}

	var signature, errSign = blockchain.signer.GenerateKeyPair()
	if errSign != nil {
		return nil, errSign
	}

	blockchain.signature = signature

	_, err := blockchain.MineBlockFromPool(creator)
	if err != nil {
		return nil, err
	}

	return &blockchain, nil
}

func (blockchain *Blockchain) createBaseTx(recipient string, fee int64) (transaction.Transaction, error) {
	var coinbaseTx, txErr = transaction.CreateTransaction(coin_transfer.CoinTransfer, EmptyAddress, int64(blockchain.CurrentRewards)+fee, 0, map[string]any{
		"recipient": recipient,
	})
	if txErr != nil {
		return nil, txErr
	}
	var signErr = coinbaseTx.AddSing(blockchain.signer, blockchain.signature)
	if signErr != nil {
		return nil, signErr
	}

	var verifyErr = coinbaseTx.Verify(blockchain.signer)
	if verifyErr != nil {
		return nil, verifyErr
	}

	return coinbaseTx, nil
}

func (blockchain *Blockchain) addBlockUnsafe(block *block.Block) error {
	if blockchain.CurrentDifficulty > block.Difficulty {
		return fmt.Errorf("block difficulty is too low")
	}
	var err = block.Verify(blockchain.signer)
	if err != nil {
		return err
	}
	for _, tx := range block.Transactions {
		err := blockchain.txProcessor.Process(tx)
		if err != nil {
			blockchain.storage.Reject()
			return err
		}
	}
	blockchain.storage.Confirm()
	blockchain.deleteExecutedTxFromPoolUnsafe(block)
	blockchain.blocks = append(blockchain.blocks, *block)
	return nil
}

func (blockchain *Blockchain) deleteExecutedTxFromPoolUnsafe(block *block.Block) {
	var newPool = []transaction.Transaction{}
	for _, tx := range blockchain.txPool {
		var checkTx = &tx
		for _, addedTx := range block.Transactions {
			if tx.GetTxId() == addedTx.GetTxId() {
				checkTx = nil
				break
			}
		}
		if checkTx != nil {
			newPool = append(newPool, *checkTx)
		}
	}
	blockchain.txPool = newPool
}

func (blockchain *Blockchain) AddBlock(block *block.Block) error {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()
	return blockchain.addBlockUnsafe(block)
}

func (blockchain *Blockchain) MineBlockFromPool(creator string) (*block.Block, error) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	var prevBlock *block.Block = nil
	if len(blockchain.blocks) > 0 {
		prevBlock = &blockchain.blocks[len(blockchain.blocks)-1]
	}
	var block, err = block.NewBlock(prevBlock, blockchain.CurrentDifficulty)
	if err != nil {
		return nil, err
	}

	var fee int64 = 0
	for _, tx := range blockchain.txPool {
		fee += tx.GetFee()
	}

	var coinbaseTx, txErr = blockchain.createBaseTx(creator, fee)
	if txErr != nil {
		return nil, txErr
	}

	block.AddTransaction(&coinbaseTx)
	for _, tx := range blockchain.txPool {
		block.AddTransaction(&tx)
	}
	var blockHash, errMine = block.Mine(0)
	if errMine != nil || [32]byte(blockHash) == [32]byte{} {
		return nil, fmt.Errorf("error while creating a block")
	}

	var errBlock = blockchain.addBlockUnsafe(block)
	if errBlock != nil {
		return nil, errBlock
	}

	return block, nil
}

func (blockchain *Blockchain) AddTransactionToPool(tx transaction.Transaction) error {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	var err = tx.Verify(blockchain.signer)
	if err != nil {
		return err
	}
	err = blockchain.txProcessor.Validate(tx)
	if err != nil {
		return err
	}

	blockchain.txPool = append(blockchain.txPool, tx)

	return nil
}

func (blockchain *Blockchain) Verify(depth int) error {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	for i := len(blockchain.blocks) - 1; i >= 0 && i >= len(blockchain.blocks)-depth; i-- {
		if i > 0 {
			if blockchain.blocks[i-1].Hash != blockchain.blocks[i].Prev {
				return fmt.Errorf("block %d: has incorrect previos hash", blockchain.blocks[i].Index)
			}
		}
		var err = blockchain.blocks[i].Verify(blockchain.signer)
		if err != nil {
			return fmt.Errorf("block %d: %s", blockchain.blocks[i].Index, err.Error())
		}
	}

	return nil
}

func (bc *Blockchain) String() string {
	var sb strings.Builder
	sb.WriteString("Blockchain{\n")
	sb.WriteString(fmt.Sprintf("  CurrentDifficult: %d bits\n", bc.CurrentDifficulty))
	sb.WriteString(fmt.Sprintf("  CurrentRewards: %d\n", bc.CurrentRewards))
	sb.WriteString(fmt.Sprintf("  TxPool: %d transactions\n", len(bc.txPool)))
	sb.WriteString(fmt.Sprintf("  Blocks: %d blocks\n", len(bc.blocks)))
	for i, blk := range bc.blocks {
		sb.WriteString(fmt.Sprintf("    Block[%d]: Index=%d, TxCount=%d\n", i, blk.Index, len(blk.Transactions)))
	}
	sb.WriteString("}")
	return sb.String()
}
