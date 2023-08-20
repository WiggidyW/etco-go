package caching

type CacheableParams interface {
	CacheKey() string
}

// type MultiCacheableParams interface {
// 	CacheKeyIndex() int
// 	CacheKeys() []string
// }

type AntiCacheableParams interface {
	AntiCacheKey() string
}

type MultiAntiCacheableParams interface {
	AntiCacheKeys() []string
}
