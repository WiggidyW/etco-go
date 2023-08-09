package naiveclient

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
)

type cachingNaiveHeadClient[P UrlParams] struct {
	client.CachingClient[
		*NaiveClientFetchParams[P],
		int32,
		cache.ExpirableData[int32],
		*naiveHeadClient[P],
	]
}

func newCachingNaiveHeadClient[P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
	minExpires time.Duration,
	bufPool *cache.BufferPool,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) cachingNaiveHeadClient[P] {
	return cachingNaiveHeadClient[P]{
		client.NewCachingClient[
			*NaiveClientFetchParams[P],
			int32,
			cache.ExpirableData[int32],
			*naiveHeadClient[P],
		](
			NewNaiveHeadClient[P](
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

type naiveHeadClient[P UrlParams] struct { // caching client
	naiveClient[int32, cache.ExpirableData[int32], P]
}

func NewNaiveHeadClient[P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
) *naiveHeadClient[P] {
	return &naiveHeadClient[P]{
		newNaiveClient[int32, cache.ExpirableData[int32], P](
			rawClient,
			useAuth,
		),
	}
}

func (nhc *naiveHeadClient[P]) Fetch(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) (*cache.ExpirableData[int32], error) {
	// fetch auth if needed
	if err := nhc.maybeFetchAuth(
		ctx,
		params,
	); err != nil {
		return nil, err
	}

	// fetch from server
	if rep, err := nhc.rawClient.FetchHead(
		ctx,
		params.urlParams.Url(),
		params.auth,
	); err != nil {
		return nil, err
	} else {
		return rep, nil
	}
}
