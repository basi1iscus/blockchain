package utils

import (
	"blockchain_demo/pkg/transaction"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

func bytes(value any) ([]byte, error) {
	switch v := value.(type) {
	case int32:
		var buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		return buf, nil
	case uint32:
		var buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, v)
		return buf, nil
	case int64:
		var buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		return buf, nil
	case uint64:
		var buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v)
		return buf, nil
	case []byte:
		return v, nil
	case transaction.HexBytes:
		return []byte(v), nil
	case string:
		return []byte(v), nil
	}

	return nil, fmt.Errorf("unsupported hashed type")
}

func GetHash(values ...any) ([]byte, error) {
	var hasher = sha256.New()
	for _, val := range values {
		b, err := bytes(val)
		if err != nil {
			return nil, err
		}
		hasher.Write(b)
	}
	hash := hasher.Sum(nil)
	return hash, nil
}

func GetHash160(prefix []byte, values ...any) ([]byte, error) {
	var hasher = ripemd160.New()
	for _, val := range values {
		b, err := bytes(val)
		if err != nil {
			return nil, err
		}
		hasher.Write(b)
	}
	hash := hasher.Sum(prefix)
	return hash, nil
}
func GetBytesFromHexParam(params map[string]any, field string) ([]byte, error) {
	value, exists := params[field]
	if !exists {
		return nil, fmt.Errorf("%s not exists in params", field)
	}
	valueStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%s is not a hex string", field)
	}
	var valueBytes, tokenErr = hex.DecodeString(valueStr)
	if tokenErr != nil {
		return nil, fmt.Errorf("unsupported %s: %s", field, valueStr)
	}
	return valueBytes, nil
}

func GetInt64FromParam(params map[string]any, field string) (uint64, error) {
	value, exists := params[field]
	if !exists {
		return 0, fmt.Errorf("%s not exists in params", field)
	}
	var valueInt uint64
	switch v := value.(type) {
	case int:
		valueInt = uint64(v)
	case int8:
		valueInt = uint64(v)
	case int16:
		valueInt = uint64(v)
	case int32:
		valueInt = uint64(v)
	case int64:
		valueInt = uint64(v)
	case uint:
		valueInt = uint64(v)
	case uint8:
		valueInt = uint64(v)
	case uint16:
		valueInt = uint64(v)
	case uint32:
		valueInt = uint64(v)
	case uint64:
		valueInt = v
	default:
		return 0, fmt.Errorf("%s is not an integer type", field)
	}
	return valueInt, nil
}

func GetEnumValueFromParam[T any](params map[string]any, field string, isValid func(s string) (T, bool)) (T, error) {
	var zeroValue T
	value, exists := params[field]
	if !exists {
		return zeroValue, fmt.Errorf("%s not exists in params", "contractType")
	}
	valueStr, ok := value.(string)
	if !ok {
		return zeroValue, fmt.Errorf("contractType must be a string")
	}
	enumValue, ok := isValid(valueStr)
	if !ok {
		return zeroValue, fmt.Errorf("contractType must be a string")
	}
	return enumValue, nil
}

func ReverseMap[T comparable, K comparable](m map[K]T) map[T]K {
	reversed := make(map[T]K)
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}