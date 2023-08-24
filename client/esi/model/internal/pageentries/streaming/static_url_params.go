package streaming

import (
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/naive"
	pe "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/pageentries"
)

type staticUrlParams struct {
	url    string
	method string
}

func (p staticUrlParams) Url() string {
	return p.url
}

func (p staticUrlParams) Method() string {
	return p.method
}

func newNaivePageParams[P pe.UrlPageParams](
	params pe.NaivePageParams[P],
	page *int,
) naive.NaiveParams[staticUrlParams] {
	return naive.NaiveParams[staticUrlParams]{
		UrlParams: staticUrlParams{
			url:    params.UrlParams.PageUrl(page),
			method: params.UrlParams.Method(),
		},
		AuthParams: params.AuthParams,
	}
}
