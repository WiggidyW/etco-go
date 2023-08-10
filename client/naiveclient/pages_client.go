package naiveclient

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
	"github.com/WiggidyW/weve-esi/util"
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

// returns a PageStream that will receive all pages.
// blocks until the number of pages is known.
func (npsc *NaivePagesClient[E, P]) FetchStreamBlocking(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) (ChanRecvPage[*client.CachingRep[[]E]], error) {
	// get num pages (blocking)
	if headRep, err := npsc.headClient.Fetch(
		ctx,
		newStaticFetchParams[P](params),
	); err != nil {
		return ChanRecvPage[*client.CachingRep[[]E]]{}, err
	} else {
		// fetch all pages in parallel (non-blocking)
		return npsc.fetchPages(ctx, params, headRep), nil
	}
}

// returns a channel that will receive a single PageStream
// non-blocking
func (npsc *NaivePagesClient[E, P]) FetchStream(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) util.ChanRecvResult[ChanRecvPage[*client.CachingRep[[]E]]] {
	// create the send and receive result channels
	chnSend, chnRecv := util.
		NewChanResult[ChanRecvPage[*client.CachingRep[[]E]]](ctx).
		Split()

	// fetch the page chan and send it on the send result chan
	go func() {
		if chnRecvPage, err := npsc.FetchStreamBlocking(
			ctx,
			params,
		); err != nil {
			chnSend.SendErr(err)
		} else {
			chnSend.SendOk(chnRecvPage)
		}
	}()

	// return the receive result chan
	return chnRecv
}

// returns all pages
func (npsc *NaivePagesClient[E, P]) FetchAll(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) (rep []*client.CachingRep[[]E], headExpires time.Time, err error) {
	// get num pages
	headRep, err := npsc.headClient.Fetch(
		ctx,
		newStaticFetchParams[P](params),
	)
	if err != nil {
		return nil, time.Time{}, err
	}

	// only one page (or an invalid number), just fetch it and return
	if headRep.Data() <= 1 {
		var pageOne int32 = 1
		if entries, err := npsc.modelClient.Fetch(
			ctx,
			newStaticFetchPageParams[P](params, &pageOne),
		); err != nil {
			return nil, time.Time{}, err
		} else {
			return []*client.CachingRep[[]E]{entries},
				headRep.Expires(),
				nil
		}
	}

	// fetch all pages in parallel
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnRecvPage := npsc.fetchPages(ctx, params, headRep)

	if pages, err := chnRecvPage.RecvAll(); err != nil {
		return nil, time.Time{}, err
	} else {
		return pages, headRep.Expires(), nil
	}
}

// returns a PageStream that will receive numPages
// non-blocking
func (npsc *NaivePagesClient[E, P]) fetchPages(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
	headRep *client.CachingRep[int32],
) ChanRecvPage[*client.CachingRep[[]E]] {
	// create the channels
	chnSendPage, chnRecvPage := NewChanPage[*client.CachingRep[[]E]](
		ctx,
		headRep.Data(),
		headRep.Expires(),
	).Split()

	// fetch the pages and send them
	for i := int32(1); i <= headRep.Data(); i++ {
		go func(params *NaiveClientFetchParams[staticUrlParams]) {
			if page, err := npsc.modelClient.Fetch(
				ctx,
				params,
			); err != nil {
				chnSendPage.SendErr(err)
			} else {
				chnSendPage.SendOk(page)
			}
		}(newStaticFetchPageParams[P](params, &i))
	}

	// return the receive channel
	return chnRecvPage
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
			key:    params.urlParams.CacheKey(),
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
			key:    params.urlParams.PageCacheKey(page),
			method: params.urlParams.Method(),
		},
		token: params.token,
		auth:  params.auth,
	}
}

func (p staticUrlParams) Url() string {
	return p.url
}
func (p staticUrlParams) CacheKey() string {
	return p.key
}
func (p staticUrlParams) Method() string {
	return p.method
}
