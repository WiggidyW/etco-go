package naiveclient

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
)

type CachingNaivePageEntriesClient[E any, P UrlParams] struct {
	client.CachingClient[
		*NaiveClientFetchParams[P],
		[]E,
		cache.ExpirableData[[]E],
		*NaivePageEntriesClient[E, P],
	]
}

func NewCachingNaivePageEntriesClient[E any, P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
	entriesPerPage int32,
	minExpires time.Duration,
	bufPool *cache.BufferPool,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) CachingNaivePageEntriesClient[E, P] {
	return CachingNaivePageEntriesClient[E, P]{
		client.NewCachingClient[
			*NaiveClientFetchParams[P],
			[]E,
			cache.ExpirableData[[]E],
			*NaivePageEntriesClient[E, P],
		](
			NewNaivePageEntriesClient[E, P](
				rawClient,
				useAuth,
				entriesPerPage,
			),
			minExpires,
			bufPool,
			clientCache,
			serverCache,
			serverLockTTL,
			serverLockMaxWait,
		),
	}
}

type NaivePageEntriesClient[E any, P UrlParams] struct {
	inner          *naivePageClient[[]E, P]
	entriesPerPage int32
}

func NewNaivePageEntriesClient[E any, P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
	entriesPerPage int32,
) *NaivePageEntriesClient[E, P] {
	return &NaivePageEntriesClient[E, P]{
		entriesPerPage: entriesPerPage,
		inner: NewNaivePageClient[[]E, P](
			rawClient,
			useAuth,
		),
	}
}

func (npec *NaivePageEntriesClient[E, P]) Fetch(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) (*cache.ExpirableData[[]E], error) {
	return npec.inner.Fetch(
		ctx,
		newNaivePageClientFetchParams[[]E, P](
			make([]E, 0, npec.entriesPerPage),
			params,
		),
	)
}

func (npec *NaivePageEntriesClient[E, P]) EntriesPerPage() int32 {
	return npec.entriesPerPage
}

type CachingNaivePageModelClient[M any, P UrlParams] struct {
	client.CachingClient[
		*NaiveClientFetchParams[P],
		M,
		cache.ExpirableData[M],
		*NaivePageModelClient[M, P],
	]
}

func NewCachingNaivePageModelClient[M any, P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
	minExpires time.Duration,
	bufPool *cache.BufferPool,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) CachingNaivePageModelClient[M, P] {
	return CachingNaivePageModelClient[M, P]{
		client.NewCachingClient[
			*NaiveClientFetchParams[P],
			M,
			cache.ExpirableData[M],
			*NaivePageModelClient[M, P],
		](
			NewNaivePageModelClient[M, P](
				rawClient,
				useAuth,
			),
			minExpires,
			bufPool,
			clientCache,
			serverCache,
			serverLockTTL,
			serverLockMaxWait,
		),
	}
}

type NaivePageModelClient[M any, P UrlParams] struct {
	inner *naivePageClient[M, P]
}

func NewNaivePageModelClient[M any, P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
) *NaivePageModelClient[M, P] {
	return &NaivePageModelClient[M, P]{
		inner: NewNaivePageClient[M, P](
			rawClient,
			useAuth,
		),
	}
}

func (npec *NaivePageModelClient[M, P]) Fetch(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) (*cache.ExpirableData[M], error) {
	var model M
	return npec.inner.Fetch(
		ctx,
		newNaivePageClientFetchParams[M, P](
			model,
			params,
		),
	)
}

type naivePageClientFetchParams[M any, P UrlParams] struct {
	Model M
	*NaiveClientFetchParams[P]
}

func newNaivePageClientFetchParams[M any, P UrlParams](
	model M,
	params *NaiveClientFetchParams[P],
) naivePageClientFetchParams[M, P] {
	return naivePageClientFetchParams[M, P]{
		Model:                  model,
		NaiveClientFetchParams: params,
	}
}

type naivePageClient[M any, P UrlParams] struct {
	naiveClient[M, cache.ExpirableData[M], P]
}

func NewNaivePageClient[M any, P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
) *naivePageClient[M, P] {
	return &naivePageClient[M, P]{
		newNaiveClient[M, cache.ExpirableData[M], P](
			rawClient,
			useAuth,
		),
	}
}

func (npc *naivePageClient[M, P]) Fetch(
	ctx context.Context,
	params naivePageClientFetchParams[M, P],
) (*cache.ExpirableData[M], error) {
	// fetch auth if needed
	if err := npc.maybeFetchAuth(
		ctx,
		params.NaiveClientFetchParams,
	); err != nil {
		return nil, err
	}

	// fetch from server
	if rep, err := rawclient.Fetch[M](
		npc.rawClient,
		ctx,
		params.urlParams.Url(),
		params.urlParams.Method(),
		params.auth,
		&params.Model,
	); err != nil {
		return nil, err
	} else {
		return rep, nil
	}
}
