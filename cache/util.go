package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func lockKey(key string) string {
	return fmt.Sprintf("%s.lock", key)
}

func deserialize[T any](b []byte) (*T, error) {
	// create an empty val
	var val T

	// create decoder
	reader := bytes.NewReader(b)
	decoder := gob.NewDecoder(reader)

	// decode bytes into &val
	err := decoder.Decode(&val)
	if err != nil {
		return nil, err
	}

	// return &val
	return &val, nil
}

func serialize[T any](val T, b *[]byte) ([]byte, error) {
	// create encoder
	buffer := bytes.NewBuffer(*b)
	encoder := gob.NewEncoder(buffer)

	// encode val into bytes
	err := encoder.Encode(val)
	if err != nil {
		return nil, err
	}

	// return bytes
	return buffer.Bytes(), nil
}
