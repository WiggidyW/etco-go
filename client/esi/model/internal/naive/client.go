package naive

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/internal/raw"
)

type NaiveClient[P UrlParams] raw.RawClient

func FetchModel[P UrlParams, M any](
	nc NaiveClient[P],
	ctx context.Context,
	params NaiveParams[P],
	model *M,
) (*cache.ExpirableData[M], error) {
	if auth, err := nc.auth(ctx, params); err != nil {
		return nil, err
	} else {
		return raw.FetchModel[M](
			raw.RawClient(nc),
			ctx,
			params.UrlParams.Url(),
			params.UrlParams.Method(),
			auth,
			model,
		)
	}
}

func (nc NaiveClient[P]) FetchHead(
	ctx context.Context,
	params NaiveParams[P],
) (*cache.ExpirableData[int], error) {
	if auth, err := nc.auth(ctx, params); err != nil {
		return nil, err
	} else {
		return raw.RawClient(nc).FetchHead(
			ctx,
			params.UrlParams.Url(),
			auth,
		)
	}
}

// returns auth token if a token is provided, else nil
// if a token is provided but auth is nil, calls fetchAuth()
func (nc NaiveClient[P]) auth(
	ctx context.Context,
	params NaiveParams[P],
) (*string, error) {
	if params.AuthParams == nil {
		return nil, nil
	} else if params.AuthParams.Auth != nil {
		return params.AuthParams.Auth, nil
	}

	if err := nc.fetchAuth(ctx, *params.AuthParams); err != nil {
		return nil, err
	}

	return params.AuthParams.Auth, nil
}

// mutates params 'auth' field using params 'token' field
func (nc NaiveClient[P]) fetchAuth(
	ctx context.Context,
	params AuthParams,
) error {
	authRep, err := raw.RawClient(nc).FetchAuth(ctx, params.Token)
	if err != nil {
		return err
	}
	params.Auth = &authRep.AccessToken
	return nil
}
