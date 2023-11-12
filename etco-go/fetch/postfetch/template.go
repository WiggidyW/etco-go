package postfetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache/expirable"
)

func CacheNamespace(
	cacheKey, typeStr string,
	expires time.Time,
) *CacheActionNamespace {
	return &CacheActionNamespace{
		CacheKey: cacheKey,
		TypeStr:  typeStr,
		Expires:  expires,
	}
}

func DualCacheSetOne[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
) []CacheActionSet {
	return []CacheActionSet{DualCacheSet(cacheKey, typeStr, data, expires)}
}
func ServerCacheSetOne[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
) []CacheActionSet {
	return []CacheActionSet{ServerCacheSet(cacheKey, typeStr, data, expires)}
}
func LocalCacheSetOne[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
) []CacheActionSet {
	return []CacheActionSet{LocalCacheSet(cacheKey, typeStr, data, expires)}
}

func DualCacheSet[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
) CacheActionSet {
	return cacheSet(cacheKey, typeStr, data, expires, true, true)
}
func ServerCacheSet[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
) CacheActionSet {
	return cacheSet(cacheKey, typeStr, data, expires, false, true)
}
func LocalCacheSet[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
) CacheActionSet {
	return cacheSet(cacheKey, typeStr, data, expires, true, false)
}

func cacheSet[T any](
	cacheKey, typeStr string,
	data *T,
	expires time.Time,
	local, server bool,
) CacheActionSet {
	return CacheActionSet{
		CacheKey:  cacheKey,
		TypeStr:   typeStr,
		Expirable: expirable.New[T](data, expires),
		Expires:   expires,
		Local:     local,
		Server:    server,
	}
}
