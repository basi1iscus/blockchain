package transaction

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/utils"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction interface {
	CalcTxId() error
	Sing(signature *sign.Signature) error
	Verify() error
}
type BaseTransaction struct {
	TxId      [32]byte
	Sender    [20]byte
	Reciver   [20]byte
	Value     int64
	Time      int64
	Sign      [65]byte
	PublicKey [91]byte
}

func NewTransaction(sender string, reciver string, value int64) (*BaseTransaction, error) {
	var senderBytes, senderErr = hex.DecodeString(sender)
	if senderErr != nil || len(senderBytes) != 20 {
		return nil, fmt.Errorf("unsupported sender format: %s", sender)
	}
	var reciverBytes, reciverErr = hex.DecodeString(reciver)
	if reciverErr != nil || len(reciverBytes) != 20 {
		return nil, fmt.Errorf("unsupported reciver format: %s", sender)
	}

	var tx = BaseTransaction{
		TxId:      [32]byte{},
		Sender:    [20]byte{},
		Reciver:   [20]byte{},
		Value:     value,
		Time:      time.Now().UnixNano(),
		Sign:      [65]byte{},
		PublicKey: [91]byte{},
	}
	copy(tx.Sender[:], senderBytes)
	copy(tx.Reciver[:], reciverBytes)

	return &tx, nil
}

func (tx *BaseTransaction) CalcTxId() error {
	var hash, err = utils.GetHash(tx.Sender[:], tx.Reciver[:], tx.Value, tx.Time)
	if err != nil {
		return err
	}
	copy(tx.TxId[:], hash)
	return nil
}

func (tx *BaseTransaction) Sing(signature *sign.Signature) error {
	var signed, err = sign.Sign(tx.TxId[:], signature.PrivateKey)
	if err != nil {
		return err
	}
	copy(tx.Sign[:], signed)
	copy(tx.PublicKey[:], signature.PublicKey)
	return nil
}

func (tx *BaseTransaction) Verify() error {
	var txToVerify = BaseTransaction{
		TxId:      [32]byte{},
		Sender:    tx.Sender,
		Reciver:   tx.Reciver,
		Value:     tx.Value,
		Time:      tx.Time,
		Sign:      [65]byte{},
		PublicKey: [91]byte{},
	}
	txToVerify.CalcTxId()
	if txToVerify.TxId != tx.TxId {
		return fmt.Errorf("TxId is invalid")
	}

	var isValid, err = sign.Verify(tx.TxId[:], tx.Sign[:], tx.PublicKey[:])
	if err != nil || !isValid {
		return fmt.Errorf("TxId signature is invalid")
	}
	return nil
}

func (tx *BaseTransaction) String() string {
	return fmt.Sprintf("Transaction{TxId: %x, Sender: %x, Reciver: %x, Value: %d, Time: %d}",
		tx.TxId, tx.Sender, tx.Reciver, tx.Value, tx.Time)
}
