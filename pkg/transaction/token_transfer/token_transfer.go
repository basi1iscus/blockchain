package token_transfer

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const TokenTransfer transaction.TransactionType = "token_transfer"

type TokenTransferTransaction struct {
	transaction.BaseTransaction
	Recipient    transaction.HexBytes `json:"recipient" json-hex:"true"`
	TokenAddress transaction.HexBytes `json:"tokenAddress" json-hex:"true"`
}

func NewTransaction(sender string, value int64, fee int64, params map[string]any) (*TokenTransferTransaction, error) {
	var senderBytes, senderErr = hex.DecodeString(sender)
	if senderErr != nil {
		return nil, fmt.Errorf("unsupported sender format: %s", sender)
	}

	var recipientBytes, recipientErr = utils.GetBytesFromHexParam(params, "recipient")
	if recipientErr != nil {
		return nil, recipientErr
	}
	var tokenBytes, tokentErr = utils.GetBytesFromHexParam(params, "token")
	if tokentErr != nil {
		return nil, tokentErr
	}

	var tx = TokenTransferTransaction{
		BaseTransaction: transaction.BaseTransaction{
			TxType:    TokenTransfer,
			TxId:      [32]byte{},
			Sender:    senderBytes,
			Value:     value,
			Fee:       fee,
			Timestamp: time.Now().UnixNano(),
			Sign:      nil,
			PublicKey: []byte{},
		},
		Recipient:    recipientBytes,
		TokenAddress: tokenBytes,
	}
	var hash, err = tx.CalcHash()
	if err != nil {
		return nil, err
	}
	tx.TxId = [32]byte(hash)
	return &tx, nil
}

func (tx *TokenTransferTransaction) GetDataForHash() []any {
	var data = tx.BaseTransaction.GetDataForHash()
	data = append(data, tx.Recipient)
	data = append(data, tx.TokenAddress)

	return data
}

func (tx *TokenTransferTransaction) CalcHash() ([]byte, error) {
	var hash, err = utils.GetHash(tx.GetDataForHash()...)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (tx *TokenTransferTransaction) Verify(signer sign.Signer) error {
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

func (tx *TokenTransferTransaction) String() string {
	return fmt.Sprintf("Transaction{TxId: %x, Sender: %x, Recipient: %x, Value: %d, Time: %d}",
		tx.TxId, tx.Sender, tx.Recipient, tx.Value, tx.Timestamp)
}

func (tx *TokenTransferTransaction) Stringify() ([]byte, error) {
	var data, err = json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func init() {
	transaction.RegisterTransactionType(TokenTransfer, func(sender string, value int64, fee int64, params map[string]any) (transaction.Transaction, error) {
		return NewTransaction(sender, value, fee, params)
	}, func() transaction.Transaction {
		return &TokenTransferTransaction{}
	})
}
