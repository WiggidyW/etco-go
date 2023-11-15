package localcache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"

	"github.com/WiggidyW/etco-go/logger"
)

func NewTypeStr[T any](desc string) (
	typeStr string,
	minBufPoolCap int,
) {
	var t T
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	logger.MaybeFatal(encoder.Encode(t))
	buf.Write([]byte(desc))
	b := buf.Bytes()
	hash := md5.Sum(b)
	typeStr = string(hash[:])
	minBufPoolCap = len(b)
	return typeStr, minBufPoolCap
}
