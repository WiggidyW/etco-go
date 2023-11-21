package esi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
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
	return fetch.FetchWithRetries[[]E](
		x,
		func(x cache.Context) ([]E, time.Time, error) {
			return getModel[[]E](
				x,
				pageUrl,
				http.MethodGet,
				auth,
				newModel,
			)
		},
		ESI_NUM_RETRIES,
		esiShouldRetry,
	)
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
	return fetch.FetchWithRetries[int](
		x,
		func(x cache.Context) (int, time.Time, error) {
			return getHead(x, url, auth)
		},
		ESI_NUM_RETRIES,
		esiShouldRetry,
	)
}
