package esi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
)

// this thing must be devastating to the compiler
func entriesNewModelFunc[E any](
	entriesPerPage int,
) func() *[]E {
	return func() *[]E {
		entries := make([]E, 0, entriesPerPage)
		return &entries
	}
}

func urlPage(url string, page int) string {
	if page == 1 {
		return url
	} else {
		return fmt.Sprintf("%s?page=%d", url, page)
	}
}

func streamGet[E any](
	x cache.Context,
	url, method string,
	entriesPerPage int,
	auth *RefreshTokenAndApp,
) (
	repOrStream RepOrStream[E],
	expires time.Time,
	pages int,
	err error,
) {
	pages, expires, err = numPagesGet(x, url, auth)
	if err != nil {
		return repOrStream, expires, 0, err
	} else if pages < 1 {
		return repOrStream, expires, 0, errors.New("ESI pages < 1")
	}

	newModel := entriesNewModelFunc[E](entriesPerPage)
	if pages == 1 {
		var repVal []E
		var pageExpires time.Time
		repVal, pageExpires, err =
			pageGet[E](x, url, method, 1, auth, newModel)
		repOrStream.Rep = &repVal
		if pageExpires.Before(expires) {
			expires = pageExpires
		}
	} else {
		chn := expirable.NewChanResult[[]E](x.Ctx(), pages, 0)
		// the last page will have fewer entries than entriesPerPage
		for page := 1; page < pages; page++ {
			go transceivePageGet(x, url, method, page, auth, newModel, chn)
		}
		go transceivePageGet(x, url, method, pages, auth, nil, chn)
		repOrStream.Stream = &chn
	}

	return repOrStream, expires, pages, err
}

func transceivePageGet[E any](
	x cache.Context,
	url, method string,
	page int,
	auth *RefreshTokenAndApp,
	newModel func() *[]E,
	chn expirable.ChanResult[[]E],
) error {
	return chn.SendExp(pageGet[E](x, url, method, page, auth, newModel))
}

func pageGet[E any](
	x cache.Context,
	url, method string,
	page int,
	auth *RefreshTokenAndApp,
	newModel func() *[]E,
) (
	rep []E,
	expires time.Time,
	err error,
) {
	pageUrl := urlPage(url, page)
	return fetch.HandleFetchVal[[]E](
		x,
		nil,
		pageGetFetchFunc[E](pageUrl, method, auth, newModel),
		EsiRetry,
	)
}

func pageGetFetchFunc[E any](
	url, method string,
	auth *RefreshTokenAndApp,
	newModel func() *[]E,
) fetch.Fetch[[]E] {
	return func(x cache.Context) (
		rep *[]E,
		expires time.Time,
		_ *postfetch.Params,
		err error,
	) {
		var repVal []E
		repVal, expires, err =
			getModel[[]E](
				x,
				url, http.MethodGet,
				auth,
				newModel,
			)
		return &repVal, expires, nil, err
	}
}

func numPagesGet(
	x cache.Context,
	url string,
	auth *RefreshTokenAndApp,
) (
	pages int,
	expires time.Time,
	err error,
) {
	var pagesPtr *int
	pagesPtr, expires, err = fetch.HandleFetch[int](
		x,
		nil,
		numPagesGetFetchFunc(url, auth),
		EsiRetry,
	)
	if pagesPtr != nil {
		pages = *pagesPtr
	}
	return pages, expires, err
}

func numPagesGetFetchFunc(
	url string,
	auth *RefreshTokenAndApp,
) fetch.Fetch[int] {
	return func(x cache.Context) (
		rep *int,
		expires time.Time,
		_ *postfetch.Params,
		err error,
	) {
		var repVal int
		repVal, expires, err = getHead(x, url, auth)
		return &repVal, expires, nil, err
	}
}
