package caching

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
)

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
