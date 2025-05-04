package transaction

import (
	"blockchain_demo/pkg/sign"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

type TransactionType string

type Hash [32]byte

func (b Hash) MarshalJSON() ([]byte, error) {
    return []byte(`"` + hex.EncodeToString(b[:]) + `"`), nil
}

func (b *Hash) UnmarshalJSON(data []byte) error {
    s, err := strconv.Unquote(string(data))
    if err != nil {
        return err
    }
    decoded, err := hex.DecodeString(s)
    if err != nil {
        return err
    }
    *b = Hash(decoded)
    return nil
}

type HexBytes []byte

func (b HexBytes) MarshalJSON() ([]byte, error) {
    return []byte(`"` + hex.EncodeToString(b[:]) + `"`), nil
}

func (b *HexBytes) UnmarshalJSON(data []byte) error {
    s, err := strconv.Unquote(string(data))
    if err != nil {
        return err
    }
    decoded, err := hex.DecodeString(s)
    if err != nil {
        return err
    }
    *b = decoded
    return nil
}

type Transaction interface {
	GetTxId() [32]byte
	GetValue() int64
	GetTime() int64
	GetSender() []byte
	AddSing(signature *sign.Signature) error
	Verify() error
	GetDataForHash() []any
	CalcHash() ([]byte, error)
	Stringify() ([]byte, error)
}

type TransactionConstructor func(sender string, value int64, fee int64, params map[string]any) (Transaction, error)

var factory = make(map[TransactionType]TransactionConstructor)
var registry = make(map[TransactionType]func() Transaction)

// Base (common) transaction struct
type BaseTransaction struct {
	TxType    TransactionType `json:"type"`
	TxId      Hash        	  `json:"id" json-hex:"true"`
	Value     int64           `json:"value"`
	Fee       int64           `json:"fee"`
	Timestamp int64           `json:"timestamp"`
	Sender    HexBytes          `json:"sender" json-hex:"true"`
	Sign      HexBytes          `json:"sign" json-hex:"true"`
	PublicKey HexBytes          `json:"public_key" json-hex:"true"`
}

func (tx *BaseTransaction) GetTxId() [32]byte {
	return tx.TxId
}
func (tx *BaseTransaction) GetValue() int64 {
	return tx.Value
}
func (tx *BaseTransaction) GetTime() int64 {
	return tx.Timestamp
}
func (tx *BaseTransaction) GetSender() []byte {
	return tx.Sender
}

func (tx *BaseTransaction) GetDataForHash() []any {
	var data = []any{}
	data = append(data, string(tx.TxType))
	data = append(data, tx.Sender)
	data = append(data, tx.Timestamp)
	data = append(data, tx.Value)
	data = append(data, tx.Fee)
	return data
}

func (tx *BaseTransaction) AddSing(signature *sign.Signature) error {
	var signed, err = sign.Sign(tx.TxId[:], signature.PrivateKey)
	if err != nil {
		return err
	}
	tx.Sign = signed
	tx.PublicKey = signature.PublicKey

	return nil
}

func (tx *BaseTransaction) Verify() error {
	var isValid, err = sign.Verify(tx.TxId[:], tx.Sign[:], tx.PublicKey[:])
	if err != nil || !isValid {
		return fmt.Errorf("TxId signature is invalid")
	}
	return nil
}

func RegisterTransactionType(name TransactionType, constructor TransactionConstructor, emptyFactory func() Transaction) {
    factory[name] = constructor
    registry[name] = emptyFactory
}

func CreateTransaction(txType TransactionType, sender string, value int64, fee int64, params map[string]any) (Transaction, error) {
    constructor, exists := factory[txType]
    if !exists {
        return nil, fmt.Errorf("transaction type %s not registered", txType)
    }

	var tx, err = constructor(sender, value, fee, params)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ParseTransaction(data []byte) (Transaction, error) {
    // Сначала определяем тип
    var typeInfo struct {
        TransactionType string `json:"type"`
    }
    
    if err := json.Unmarshal(data, &typeInfo); err != nil {
        return nil, fmt.Errorf("failed to parse transaction type: %v", err)
    }
    
    // Создаем экземпляр нужного типа
    emptyFactory, exists := registry[TransactionType(typeInfo.TransactionType)]
    if !exists {
        return nil, fmt.Errorf("unknown transaction type: %s", typeInfo.TransactionType)
    }
    
    tx := emptyFactory()
    
    // Десериализуем полные данные
    if err := json.Unmarshal(data, tx); err != nil {
        return nil, fmt.Errorf("failed to parse transaction data: %v", err)
    }
    
    return tx, nil
}

func (tx *BaseTransaction) Stringify() ([]byte, error) {
    var data, err = json.Marshal(tx) 
	if err != nil {
        return nil, fmt.Errorf("failed to serialize transaction: %v", err)
    }
    
    return data, nil
}