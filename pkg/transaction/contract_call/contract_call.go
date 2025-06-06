package contract_call

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const ContractCall transaction.TransactionType = "contract_call"

type ContractCallParams struct {
	To     transaction.HexBytes `json:"to" json-hex:"true"`
	Amount uint64 				`json:"amount"`
}

type ContractCallTransaction struct {
	transaction.BaseTransaction
	ContractAddress transaction.HexBytes       `json:"contractAddress" json-hex:"true"`
	Method          ContractMethod 			   `json:"method"`
	InitParams      ContractCallParams
}

type ContractMethod string

const (
	Transfer ContractMethod = "transfer"
)

// Или с использованием карты для более эффективной проверки
var validContractTypes = map[ContractMethod]bool{
	Transfer: true,
}

func IsValid(s string) (ContractMethod, bool) {
	return ContractMethod(s), validContractTypes[ContractMethod(s)]
}

func NewTransaction(sender string, value int64, fee int64, params map[string]any) (*ContractCallTransaction, error) {
	var senderBytes, senderErr = hex.DecodeString(sender)
	if senderErr != nil || len(senderBytes) != 20 {
		return nil, fmt.Errorf("unsupported sender format: %s", sender)
	}

	var contractAddressBytes, contractAddressErr = utils.GetBytesFromHexParam(params, "contractAddress")
	if contractAddressErr != nil {
		return nil, contractAddressErr
	}
	var toBytes, toErr = utils.GetBytesFromHexParam(params, "to")
	if toErr != nil {
		return nil, toErr
	}
	var amount, amountErr = utils.GetInt64FromParam(params, "amount")
	if amountErr != nil {
		return nil, amountErr
	}
	var method, methodErr = utils.GetEnumValueFromParam(params, "method", IsValid)
	if methodErr != nil {
		return nil, methodErr
	}

	var tx = ContractCallTransaction{
		BaseTransaction: transaction.BaseTransaction{
			TxType:    ContractCall,
			TxId:      [32]byte{},
			Sender:    senderBytes,
			Value:     value,
			Fee:       fee,
			Timestamp: time.Now().UnixNano(),
			Sign:      nil,
			PublicKey: []byte{},
		},
		ContractAddress: contractAddressBytes,
		Method:    method,
		InitParams:      ContractCallParams{To: toBytes, Amount: amount},
	}
	var hash, err = tx.CalcHash()
	if err != nil {
		return nil, err
	}
	tx.TxId = [32]byte(hash)
	return &tx, nil
}

func (tx *ContractCallTransaction) GetDataForHash() []any {
	var data = tx.BaseTransaction.GetDataForHash()
	data = append(data, tx.ContractAddress)
	data = append(data, string(tx.Method))
	data = append(data, tx.InitParams.To)
	data = append(data, tx.InitParams.Amount)

	return data
}

func (tx *ContractCallTransaction) CalcHash() ([]byte, error) {
	var hash, err = utils.GetHash(tx.GetDataForHash()...)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (tx *ContractCallTransaction) Verify(signer sign.Signer) error {
	var hash, hashErr = tx.CalcHash()
	if hashErr != nil {
		return fmt.Errorf("unable to calculate hash")
	}
	if [32]byte(hash) != tx.TxId {
		return fmt.Errorf("TxId is invalid")
	}
	var err = tx.BaseTransaction.Verify(signer)
	if err != nil {
		return fmt.Errorf("TxId signature is invalid")
	}
	return nil
}

func (tx *ContractCallTransaction) String() string {
	return fmt.Sprintf("Transaction{TxId: %x, Sender: %x, Contract: %x, Value: %d, Time: %d}",
		tx.TxId, tx.Sender, tx.ContractAddress, tx.Value, tx.Timestamp)
}

func (tx *ContractCallTransaction) Stringify() ([]byte, error) {
	var data, err = json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func init() {
	transaction.RegisterTransactionType(ContractCall, func(sender string, value int64, fee int64, params map[string]any) (transaction.Transaction, error) {
		return NewTransaction(sender, value, fee, params)
	}, func() transaction.Transaction {
		return &ContractCallTransaction{}
	})
}
