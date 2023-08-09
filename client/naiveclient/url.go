package naiveclient

type UrlParams interface {
	Url() string
	Key() string
	Method() string
}

type UrlPageParams interface {
	UrlParams
	PageUrl(page *int32) string
	PageKey(page *int32) string
}
