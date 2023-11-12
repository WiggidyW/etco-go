package esi

import (
	"errors"
	"net/http"
	"time"

	built "github.com/WiggidyW/etco-go/builtinconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/esierror"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

var (
	ForbiddenStructure *StructureInfo = &StructureInfo{
		Forbidden:     true,
		Name:          built.FORBIDDEN_STRUCTURE_NAME,
		SolarSystemId: built.FORBIDDEN_SYSTEM_ID,
	}
)

func infoGet[M any](
	x cache.Context,
	url, method, cacheKey, typeStr string,
	lockTTL, lockMaxBackoff, minExpiresIn time.Duration,
	auth *RefreshTokenAndApp,
	handleErr func(error) (ok bool, rep *M, expires time.Time),
) (
	rep *M,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch[M](
		x,
		&prefetch.Params[M]{
			CacheParams: &prefetch.CacheParams[M]{
				Get: prefetch.DualCacheGet[M](
					cacheKey, typeStr,
					true,
					nil,
					cache.SloshTrue[M],
				),
			},
		},
		infoGetFetchFunc[M](
			url, method, cacheKey, typeStr,
			minExpiresIn,
			auth,
			handleErr,
		),
		EsiRetry,
	)
}

func infoGetFetchFunc[M any](
	url, method, cacheKey, typeStr string,
	minExpiresIn time.Duration,
	auth *RefreshTokenAndApp,
	handleErr func(error) (ok bool, rep *M, expires time.Time),
) fetch.Fetch[M] {
	return func(x cache.Context) (
		rep *M,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var repVal M
		repVal, expires, err = getModel[M](x, url, method, auth, nil)
		if err != nil {
			var ok bool
			ok, rep, expires = handleErr(err)
			if !ok {
				return nil, expires, nil, err
			}
		} else {
			rep = &repVal
		}
		expires = fetch.CalcExpires(expires, minExpiresIn)
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.DualCacheSetOne(
					cacheKey, typeStr,
					rep,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}

func entityInfoHandleErr[M any](err error) (
	ok bool,
	rep *M,
	expires time.Time,
) {
	var statusErr esierror.StatusError
	ok = errors.As(err, &statusErr) && statusErr.Code == http.StatusNotFound
	return ok, nil, expires
}

func structureInfoHandleErr(err error) (
	ok bool,
	rep *StructureInfo,
	expires time.Time,
) {
	var statusErr esierror.StatusError
	if errors.As(err, &statusErr) {
		if statusErr.Code == http.StatusNotFound {
			ok = true
		} else if statusErr.Code == http.StatusForbidden {
			ok = true
			rep = ForbiddenStructure
		}
	}
	return ok, rep, expires
}
