package contract_deploy_processor

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/transaction/contract_deploy"
	"blockchain_demo/pkg/transaction_processor"
	"fmt"
)

type ContractDeployProcessor struct {
	storage ballance_storage.BallanceStorage
}

func NewProcessor(storage ballance_storage.BallanceStorage) transaction_processor.TransactionProcessor {
	var processor = ContractDeployProcessor{
		storage: storage,
	}

	return &processor
}

func (p *ContractDeployProcessor) Validate(tx transaction.Transaction) error {

	coinTx, ok := tx.(*contract_deploy.ContractDeployTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	ballance := p.storage.GetBallance(string(coinTx.Sender))

	if ballance < int64(coinTx.InitParams.InitialSupplay) {
		return fmt.Errorf("sender's balance is too low")
	}
	return nil
}

func (p *ContractDeployProcessor) Process(tx transaction.Transaction) error {
	coinTx, ok := tx.(*contract_deploy.ContractDeployTransaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := p.Validate(tx); err != nil {
		return err
	}

	total := int64(coinTx.InitParams.InitialSupplay) + tx.GetFee()
	p.storage.SubBallance(string(coinTx.Sender), total)
	p.storage.AddBallance(string(coinTx.ContractAddress), int64(coinTx.InitParams.InitialSupplay))

	return nil
}
