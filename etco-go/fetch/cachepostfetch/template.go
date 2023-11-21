package cachepostfetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache/expirable"
)

func Namespace(
	cacheKey, typeStr string,
	expires time.Time,
) *ActionNamespace {
	return &ActionNamespace{
		CacheKey: cacheKey,
		TypeStr:  typeStr,
		Expires:  expires,
	}
}

func DualSetOne[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
) []ActionSet {
	return []ActionSet{DualSet(cacheKey, typeStr, data, expires)}
}
func ServerSetOne[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
) []ActionSet {
	return []ActionSet{ServerSet(cacheKey, typeStr, data, expires)}
}
func LocalSetOne[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
) []ActionSet {
	return []ActionSet{LocalSet(cacheKey, typeStr, data, expires)}
}

func DualSet[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
) ActionSet {
	return cacheSet(cacheKey, typeStr, data, expires, true, true)
}
func ServerSet[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
) ActionSet {
	return cacheSet(cacheKey, typeStr, data, expires, false, true)
}
func LocalSet[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
) ActionSet {
	return cacheSet(cacheKey, typeStr, data, expires, true, false)
}

func cacheSet[T any](
	cacheKey, typeStr string,
	data T,
	expires time.Time,
	local, server bool,
) ActionSet {
	return ActionSet{
		CacheKey:  cacheKey,
		TypeStr:   typeStr,
		Expirable: expirable.New[T](data, expires),
		Expires:   expires,
		Local:     local,
		Server:    server,
	}
}
