package postfetch

func DualCacheSet(
	typeStr string,
	cacheKey string,
) *CacheParams {
	return &CacheParams{
		TypeStr:   typeStr,
		Key:       cacheKey,
		LocalSet:  true,
		ServerSet: true,
	}
}

func ServerCacheSet(
	typeStr string,
	cacheKey string,
) *CacheParams {
	return &CacheParams{
		TypeStr:   typeStr,
		Key:       cacheKey,
		LocalSet:  false,
		ServerSet: true,
	}
}

func LocalCacheSet(
	typeStr string,
	cacheKey string,
) *CacheParams {
	return &CacheParams{
		TypeStr:   typeStr,
		Key:       cacheKey,
		LocalSet:  true,
		ServerSet: false,
	}
}
