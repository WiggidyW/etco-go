package internal

import (
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/logger"
)

func validateKeyLength(keys []string, expectLen int) {
	if len(keys) != expectLen {
		logger.Fatal(fmt.Errorf(
			"expected %d CacheKeys, got %d",
			expectLen,
			len(keys),
		))
	}
}

func NewMinExpirableData[D any, ED cache.Expirable[D]](
	rep ED,
	minExpires time.Duration,
) cache.ExpirableData[D] {
	repExpiresTime := rep.Expires()
	minExpiresTime := time.Now().Add(minExpires)
	if repExpiresTime.Before(minExpiresTime) {
		repExpiresTime = minExpiresTime
	}
	return cache.NewExpirableData[D](
		rep.Data(),
		repExpiresTime,
	)
}
