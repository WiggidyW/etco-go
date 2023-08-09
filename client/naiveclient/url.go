package naiveclient

type UrlParams interface {
	Url() string
	CacheKey() string
	Method() string
}

type UrlPageParams interface {
	UrlParams
	PageUrl(page *int32) string
	PageCacheKey(page *int32) string
}
