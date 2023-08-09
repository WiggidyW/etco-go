package naiveclient

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
)

type NaiveClientFetchParams[P UrlParams] struct {
	urlParams P
	token     *string
	auth      *string
}

func NewNaiveClientFetchParams[P UrlParams](
	urlParams P,
	token *string,
	auth *string,
) *NaiveClientFetchParams[P] {
	return &NaiveClientFetchParams[P]{
		urlParams: urlParams,
		token:     token,
		auth:      auth,
	}
}

func (ncfp NaiveClientFetchParams[P]) CacheKey() string {
	return ncfp.urlParams.CacheKey()
}

type naiveClient[D any, ED cache.Expirable[D], P UrlParams] struct {
	rawClient *rawclient.RawClient
	useAuth   bool
}

func newNaiveClient[D any, ED cache.Expirable[D], P UrlParams](
	rawClient *rawclient.RawClient,
	useAuth bool,
) naiveClient[D, ED, P] {
	return naiveClient[D, ED, P]{
		rawClient: rawClient,
		useAuth:   useAuth,
	}
}

func (nc *naiveClient[D, ED, P]) maybeFetchAuth(
	ctx context.Context,
	params *NaiveClientFetchParams[P],
) error {
	// check if fetching auth is needed
	if !nc.useAuth || params.auth != nil {
		return nil
	} else if params.token == nil {
		panic("maybeFetchAuth: useAuth == true, token == nil")
	}

	// fetch the auth
	authRep, err := nc.rawClient.FetchAuth(ctx, *params.token)
	if err != nil {
		return err
	}
	params.auth = &authRep.AccessToken

	return nil
}
