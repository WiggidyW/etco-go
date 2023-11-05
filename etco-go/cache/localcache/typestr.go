package localcache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"

	"github.com/WiggidyW/etco-go/logger"
)

func NewTypeStr[T any]() string {
	var t T
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	logger.MaybeFatal(encoder.Encode(t))
	b := buf.Bytes()
	hash := md5.Sum(b)
	typeStr := string(hash[:])
	return typeStr
}
