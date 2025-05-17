package token_transfer_processor

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/token_transfer"
	"blockchain_demo/pkg/transaction_processor"
	"fmt"
)

type TokenTransferProcessor struct {
	storage ballance_storage.BallanceStorage
}

func NewProcessor(storage ballance_storage.BallanceStorage) transaction_processor.TransactionProcessor {
	var processor = TokenTransferProcessor{
		storage: storage,
	}

	return &processor
}

func (p *TokenTransferProcessor) Validate(tx transaction.Transaction) error {

	coinTx, ok := tx.(*token_transfer.TokenTransferTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	ballance := p.storage.GetBallance(string(coinTx.Sender))

	if ballance < tx.GetValue()+tx.GetFee() {
		return fmt.Errorf("sender's balance is too low")
	}
	return nil
}

func (p *TokenTransferProcessor) Process(tx transaction.Transaction) error {
	coinTx, ok := tx.(*token_transfer.TokenTransferTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := p.Validate(tx); err != nil {
		return err
	}

	p.storage.SubBallance(string(coinTx.Sender), tx.GetValue()+tx.GetFee())
	p.storage.AddBallance(string(coinTx.Recipient), tx.GetValue())

	return nil
}
