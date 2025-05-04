package contract_deploy

import (
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const ContractDeploy transaction.TransactionType = "contract_deploy"

type InitParams struct {
	Owner          []byte `json:"owner" json-hex:"true"`
	InitialSupplay uint64 `json:"initialSupply"`
}

type ContractDeployTransaction struct {
	transaction.BaseTransaction
	Code       []byte `json:"code" json-hex:"true"`
	InitParams InitParams
}

func NewTransaction(sender string, value int64, fee int64, params map[string]any) (*ContractDeployTransaction, error) {
	var senderBytes, senderErr = hex.DecodeString(sender)
	if senderErr != nil || len(senderBytes) != 20 {
		return nil, fmt.Errorf("unsupported sender format: %s", sender)
	}
	var codeBytes, codeErr = utils.GetBytesFromHexParam(params, "code")
	if codeErr != nil {
		return nil, codeErr
	}
	var ownerBytes, ownerErr = utils.GetBytesFromHexParam(params, "owner")
	if ownerErr != nil {
		return nil, ownerErr
	}
	var initialSupplay, initialSupplayErr = utils.GetInt64FromParam(params, "initialSupplay")
	if initialSupplayErr != nil {
		return nil, initialSupplayErr
	}

	var tx = ContractDeployTransaction{
		BaseTransaction: transaction.BaseTransaction{
			TxType:    ContractDeploy,
			TxId:      [32]byte{},
			Sender:    senderBytes,
			Value:     value,
			Fee:       fee,
			Timestamp: time.Now().UnixNano(),
			Sign:      nil,
			PublicKey: []byte{},
		},
		Code: codeBytes,
		InitParams: InitParams{
			Owner:          ownerBytes,
			InitialSupplay: initialSupplay,
		},
	}
	var hash, err = tx.CalcHash()
	if err != nil {
		return nil, err
	}
	tx.TxId = [32]byte(hash)
	return &tx, nil
}

func (tx *ContractDeployTransaction) getDataForHash() []any {
	var data = tx.BaseTransaction.GetDataForHash()
	data = append(data, tx.Code)
	data = append(data, tx.InitParams.InitialSupplay)
	data = append(data, tx.InitParams.Owner)
	return data
}

func (tx *ContractDeployTransaction) CalcHash() ([]byte, error) {
	var hash, err = utils.GetHash(tx.getDataForHash()...)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (tx *ContractDeployTransaction) Verify() error {
	var hash, hashErr = tx.CalcHash()
	if hashErr != nil {
		return fmt.Errorf("unable to calculate hash")
	}
	if [32]byte(hash) != tx.TxId {
		return fmt.Errorf("TxId is invalid")
	}
	var err = tx.BaseTransaction.Verify()
	if err != nil {
		return fmt.Errorf("TxId signature is invalid")
	}
	return nil
}

func (tx *ContractDeployTransaction) String() string {
	return fmt.Sprintf("Transaction{TxId: %x, Sender: %x, Owner: %x, Value: %d, Time: %d}",
		tx.TxId, tx.Sender, tx.InitParams.Owner, tx.Value, tx.Timestamp)
}

func (tx *ContractDeployTransaction) Stringify() ([]byte, error) {
	var data, err = json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func init() {
	transaction.RegisterTransactionType(ContractDeploy, func(sender string, value int64, fee int64, params map[string]any) (transaction.Transaction, error) {
		return NewTransaction(sender, value, fee, params)
	}, func() transaction.Transaction {
		return &ContractDeployTransaction{}
	})
}
