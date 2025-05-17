package contract_call_processor

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/contract_call"
	"blockchain_demo/pkg/transaction_processor"
	"fmt"
)

type ContractCallProcessor struct {
	storage ballance_storage.BallanceStorage
}

func NewProcessor(storage ballance_storage.BallanceStorage) transaction_processor.TransactionProcessor {
	var processor = ContractCallProcessor{
		storage: storage,
	}

	return &processor
}

func (p *ContractCallProcessor) Validate(tx transaction.Transaction) error {

	coinTx, ok := tx.(*contract_call.ContractCallTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	ballance := p.storage.GetBallance(string(coinTx.ContractAddress))

	if ballance < int64(coinTx.InitParams.Amount)+coinTx.Fee {
		return fmt.Errorf("sender's balance is too low")
	}
	return nil
}

func (p *ContractCallProcessor) Process(tx transaction.Transaction) error {
	coinTx, ok := tx.(*contract_call.ContractCallTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := p.Validate(tx); err != nil {
		return err
	}

	total := int64(coinTx.InitParams.Amount) + tx.GetFee()
	p.storage.SubBallance(string(coinTx.Sender), total)
	p.storage.AddBallance(string(coinTx.InitParams.To), int64(coinTx.InitParams.Amount))

	return nil
}
