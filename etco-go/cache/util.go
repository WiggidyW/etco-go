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

// func deserialize[T any](
// 	serializedVal []byte,
// 	targetVal *T,
// ) (
// 	err error,
// ) {
// 	reader := bytes.NewReader(serializedVal)
// 	decoder := gob.NewDecoder(reader)
// 	return decoder.Decode(targetVal)
// }

// func serialize[T any](
// 	val T,
// 	buf *[]byte,
// ) (
// 	serializedVal []byte,
// 	err error,
// ) {
// 	// create encoder
// 	buffer := bytes.NewBuffer(*buf)
// 	encoder := gob.NewEncoder(buffer)

// 	// encode val into bytes
// 	err = encoder.Encode(val)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// return bytes
// 	return buffer.Bytes(), nil
// }
