package pageentries

import "github.com/WiggidyW/etco-go/client/esi/model/internal/naive"

type NaivePageParams[P UrlPageParams] struct {
	UrlParams  P
	AuthParams *naive.AuthParams
}

type UrlPageParams interface {
	PageUrl(page *int) string
	Method() string
}
