package coin_transfer_processor

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/coin_transfer"
	"blockchain_demo/pkg/transaction_processor"
	"fmt"
)

type CoinTransferProcessor struct {
	storage ballance_storage.BallanceStorage
}

func NewProcessor(storage ballance_storage.BallanceStorage) transaction_processor.TransactionProcessor {
	var processor = CoinTransferProcessor{
		storage: storage,
	}

	return &processor
}

func (p *CoinTransferProcessor) Validate(tx transaction.Transaction) error {

	coinTx, ok := tx.(*coin_transfer.CoinTransferTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if string(coinTx.Sender[:]) == string(coin_transfer.EmptyAddress[:]) {
		return nil
	}
	ballance := p.storage.GetBallance(string(coinTx.Sender))

	if ballance < tx.GetValue()+tx.GetFee() {
		return fmt.Errorf("sender's balance is too low")
	}
	return nil
}

func (p *CoinTransferProcessor) Process(tx transaction.Transaction) error {
	coinTx, ok := tx.(*coin_transfer.CoinTransferTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := p.Validate(tx); err != nil {
		return err
	}

	if string(coinTx.Sender[:]) != string(coin_transfer.EmptyAddress[:]) {
		p.storage.SubBallance(string(coinTx.Sender), tx.GetValue()+tx.GetFee())
	}

	p.storage.AddBallance(string(coinTx.Recipient), tx.GetValue())

	return nil
}
