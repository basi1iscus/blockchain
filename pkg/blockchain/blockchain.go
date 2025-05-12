package blockchain

import (
	"blockchain_demo/pkg/block"
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"fmt"
	"strings"
)

const EmptyAddress = "0000000000000000000000000000000000000000"

type Blockchain struct {
	CurrentDifficult uint64
	CurrentRewards   uint64
	TxPool           []transaction.Transaction
	Blocks           []block.Block
	signature        *sign.SignatureKeys
	signer			 sign.Signer
}

func NewBlockchain(rewards uint64, difficulty uint64, creator string, signer sign.Signer) (*Blockchain, error) {
	var blockchain = Blockchain{
		CurrentDifficult: difficulty,
		CurrentRewards:   rewards,
		TxPool:           []transaction.Transaction{},
		Blocks:           []block.Block{},
		signer: 		  signer,
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

func (blockchain *Blockchain) createBaseTx(recipient string) (transaction.Transaction, error) {
	var coinbaseTx, txErr = transaction.CreateTransaction(coin_transfer.CoinTransfer, EmptyAddress, int64(blockchain.CurrentRewards), 0, map[string]any{
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

func (blockchain *Blockchain) MineBlockFromPool(creator string) (*block.Block, error) {
	var prevBlock *block.Block = nil
	if len(blockchain.Blocks) > 0 {
		prevBlock = &blockchain.Blocks[len(blockchain.Blocks)-1]
	}
	var block, err = block.NewBlock(prevBlock, blockchain.CurrentDifficult)
	if err != nil {
		return nil, err
	}

	var coinbaseTx, txErr = blockchain.createBaseTx(creator)
	if txErr != nil {
		return nil, txErr
	}

	block.AddTransaction(&coinbaseTx)
	for _, tx := range blockchain.TxPool {
		block.AddTransaction(&tx)
	}
	var blockHash, errMine = block.Mine(0)
	if errMine != nil || [32]byte(blockHash) == [32]byte{} {
		return nil, fmt.Errorf("error while creating a block")
	}

	var errBlock = blockchain.AddBlock(block)
	if errBlock != nil {
		return nil, errBlock
	}

	return block, nil
}

func (blockchain *Blockchain) deleteExecutedTxFromPool(block *block.Block) {
	var newPool = []transaction.Transaction{}
	for _, tx := range blockchain.TxPool {
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
	blockchain.TxPool = newPool
}

func (blockchain *Blockchain) AddTransactionToPool(tx transaction.Transaction) error {
	var err = tx.Verify(blockchain.signer)
	if err != nil {
		return err
	}
	blockchain.TxPool = append(blockchain.TxPool, tx)

	return nil
}

func (blockchain *Blockchain) AddBlock(block *block.Block) error {
	var err = block.Verify(blockchain.signer)
	if err != nil {
		return err
	}
	blockchain.deleteExecutedTxFromPool(block)
	blockchain.Blocks = append(blockchain.Blocks, *block)

	return nil
}

func (blockchain *Blockchain) Verify(depth int) error {
	for i := len(blockchain.Blocks) - 1; i >= 0 && i >= len(blockchain.Blocks)-depth; i-- {
		if i > 0 {
			if blockchain.Blocks[i-1].Hash != blockchain.Blocks[i].Prev {
				return fmt.Errorf("block %d: has incorrect previos hash", blockchain.Blocks[i].Index)
			}
		}
		var err = blockchain.Blocks[i].Verify(blockchain.signer)
		if err != nil {
			return fmt.Errorf("block %d: %s", blockchain.Blocks[i].Index, err.Error())
		}
	}

	return nil
}

func (bc *Blockchain) String() string {
	var sb strings.Builder
	sb.WriteString("Blockchain{\n")
	sb.WriteString(fmt.Sprintf("  CurrentDifficult: %d bits\n", bc.CurrentDifficult))
	sb.WriteString(fmt.Sprintf("  CurrentRewards: %d\n", bc.CurrentRewards))
	sb.WriteString(fmt.Sprintf("  TxPool: %d transactions\n", len(bc.TxPool)))
	sb.WriteString(fmt.Sprintf("  Blocks: %d blocks\n", len(bc.Blocks)))
	for i, blk := range bc.Blocks {
		sb.WriteString(fmt.Sprintf("    Block[%d]: Index=%d, TxCount=%d\n", i, blk.Index, len(blk.Transactions)))
	}
	sb.WriteString("}")
	return sb.String()
}
