package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
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
