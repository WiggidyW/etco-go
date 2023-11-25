package localcache

import (
	"bytes"
	"encoding/gob"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/logger"
)

func NewTypeStr[T any](desc string) (
	typeStr keys.Key,
	minBufPoolCap int,
) {
	var t T
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	logger.MaybeFatal(encoder.Encode(t))
	b := buf.Bytes()
	minBufPoolCap = len(b)
	typeStr = keys.NewTypeStr(b, desc)
	return typeStr, minBufPoolCap
}
