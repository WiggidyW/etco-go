package naiveclient

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
)

type NaivePagesClient[E any, P UrlPageParams] struct { // non-caching client
	modelClient *CachingNaivePageEntriesClient[E, staticUrlParams]
	headClient  *cachingNaiveHeadClient[staticUrlParams]
}

func NewNaivePagesClient[E any, P UrlPageParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
	entriesPerPage int32,
	minExpires time.Duration,
	modelPool *cache.BufferPool,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	headPool *cache.BufferPool,
	headServerLockTTL time.Duration,
	headServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) NaivePagesClient[E, P] {
	modelClient := NewCachingNaivePageEntriesClient[E, staticUrlParams](
		rawClient,
		useAuth,
		entriesPerPage,
		minExpires,
		headPool,
		clientCache,
		serverCache,
		headServerLockTTL,
		headServerLockMaxWait,
	)
	headClient := newCachingNaiveHeadClient[staticUrlParams](
		rawClient,
		useAuth,
		minExpires,
		headPool,
		clientCache,
		serverCache,
		headServerLockTTL,
		headServerLockMaxWait,
	)
	return NaivePagesClient[E, P]{&modelClient, &headClient}
}

// returns a PageStream that will receive numPages
// non-blocking
func (npsc *NaivePagesClient[E, P]) fetchPages(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
	numPages *client.CachingRep[int32],
) PageStream[client.CachingRep[[]E]] {
	strm := npsc.makePageStream(numPages.Data(), numPages.Expires())
	var i int32
	for i = 1; i <= numPages.Data(); i++ {
		go func(params *NaiveClientFetchParams[staticUrlParams]) {
			if page, err := npsc.modelClient.Fetch(
				ctx,
				params,
			); err != nil {
				strm.sendErr(err)
			} else {
				strm.sendOk(page)
			}
		}(newStaticFetchPageParams[P](params, &i))
	}
	return strm
}

// returns a channel that will receive a single PageStream
// non-blocking
func (npsc *NaivePagesClient[E, P]) FetchStream(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) <-chan PageStream[client.CachingRep[[]E]] {
	chn := npsc.makePageStreamChan(1)
	go func() {
		strm := npsc.FetchStreamBlocking(ctx, params)
		chn <- strm
	}()
	return chn
}

// returns a PageStream that will receive all pages.
// blocks until the number of pages is known.
func (npsc *NaivePagesClient[E, P]) FetchStreamBlocking(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) PageStream[client.CachingRep[[]E]] {
	// get num pages (blocking)
	numPages, err := npsc.headClient.Fetch(
		ctx,
		newStaticFetchParams[P](params),
	)
	if err != nil {
		strm := npsc.makePageStream(1, time.Time{})
		strm.sendErr(err)
		return strm
	}

	// fetch all pages in parallel (non-blocking)
	return npsc.fetchPages(ctx, params, numPages)
}

// returns all pages
func (npsc *NaivePagesClient[E, P]) FetchAll(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) ([]*client.CachingRep[[]E], time.Time, error) {
	// get num pages
	numPages, err := npsc.headClient.Fetch(
		ctx,
		newStaticFetchParams[P](params),
	)
	if err != nil {
		return nil, time.Time{}, err
	}

	// only one page (or an invalid number), just fetch it and return
	if numPages.Data() <= 1 {
		var pageOne int32 = 1
		if entries, err := npsc.modelClient.Fetch(
			ctx,
			newStaticFetchPageParams[P](params, &pageOne),
		); err != nil {
			return nil, time.Time{}, err
		} else {
			return npsc.makePagesSingle(entries),
				entries.Expires(),
				nil
		}
	}

	// fetch all pages in parallel
	strm := npsc.fetchPages(ctx, params, numPages)
	defer strm.Close()

	// collect the results
	pages := npsc.makePages(numPages.Data())
	remaining := numPages.Data()
	for remaining > 0 {
		if page, err := strm.Recv(); err != nil {
			return nil, time.Time{}, err
		} else {
			pages = append(pages, page)
			remaining--
		}
	}

	return pages, numPages.Expires(), nil
}

func (NaivePagesClient[E, P]) makePages(
	numPages int32,
) []*client.CachingRep[[]E] {
	return make(
		[]*client.CachingRep[[]E],
		0,
		numPages,
	)
}

func (NaivePagesClient[E, P]) makePagesSingle(
	page *client.CachingRep[[]E],
) []*client.CachingRep[[]E] {
	return []*client.CachingRep[[]E]{page}
}

func (NaivePagesClient[E, P]) makePageStream(
	numPages int32,
	headExpires time.Time,
) PageStream[client.CachingRep[[]E]] {
	return makePageStream[client.CachingRep[[]E]](
		numPages,
		headExpires,
	)
}

func (NaivePagesClient[E, P]) makePageStreamChan(
	capacity int,
) chan PageStream[client.CachingRep[[]E]] {
	return make(
		chan PageStream[client.CachingRep[[]E]],
		capacity,
	)
}

func (npsc *NaivePagesClient[E, P]) EntriesPerPage() int32 {
	return npsc.modelClient.Client.EntriesPerPage()
}

type staticUrlParams struct {
	url    string
	key    string
	method string
}

func newStaticFetchParams[P UrlParams](
	params *NaiveClientFetchParams[P],
) *NaiveClientFetchParams[staticUrlParams] {
	return &NaiveClientFetchParams[staticUrlParams]{
		urlParams: staticUrlParams{
			url:    params.urlParams.Url(),
			key:    params.urlParams.Key(),
			method: params.urlParams.Method(),
		},
		token: params.token,
		auth:  params.auth,
	}
}
func newStaticFetchPageParams[P UrlPageParams](
	params *NaiveClientFetchParams[P],
	page *int32,
) *NaiveClientFetchParams[staticUrlParams] {
	return &NaiveClientFetchParams[staticUrlParams]{
		urlParams: staticUrlParams{
			url:    params.urlParams.PageUrl(page),
			key:    params.urlParams.PageKey(page),
			method: params.urlParams.Method(),
		},
		token: params.token,
		auth:  params.auth,
	}
}

func (p staticUrlParams) Url() string {
	return p.url
}
func (p staticUrlParams) Key() string {
	return p.key
}
func (p staticUrlParams) Method() string {
	return p.method
}
