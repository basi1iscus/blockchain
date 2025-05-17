package transaction_processor

import (
	"blockchain_demo/pkg/ballance_storage"
	"blockchain_demo/pkg/transaction"
	"fmt"
)

type TransactionProcessor interface {
	Validate(tx transaction.Transaction) error
	Process(tx transaction.Transaction) error
}

type BaseProcessor struct {
}

type ProcessorConstructor func(storage ballance_storage.BallanceStorage) TransactionProcessor

var processors = make(map[transaction.TransactionType]TransactionProcessor)

func RegisterProcessorType(name transaction.TransactionType, constructor TransactionProcessor) {
	processors[name] = constructor
}

func (v BaseProcessor) Validate(tx transaction.Transaction) error {
	processor, exist := processors[tx.GetTxType()]
	if !exist {
		return fmt.Errorf("not registered transaction processor for type: %s", tx.GetTxType())
	}

	return processor.Validate(tx)
}

func (v BaseProcessor) Process(tx transaction.Transaction) error {
	processor, exist := processors[tx.GetTxType()]
	if !exist {
		return fmt.Errorf("not registered transaction processor for type: %s", tx.GetTxType())
	}

	return processor.Process(tx)
}
