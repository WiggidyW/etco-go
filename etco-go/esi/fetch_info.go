package esi

import (
	"errors"
	"net/http"
	"time"

	built "github.com/WiggidyW/etco-go/builtinconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/esierror"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
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
	minExpiresIn time.Duration,
	auth *RefreshTokenAndApp,
	handleErr func(error) (ok bool, rep *M, expires time.Time),
) (
	rep *M,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache[*M](
		x,
		infoGetFetchFunc[M](
			url, method, cacheKey, typeStr,
			minExpiresIn,
			auth,
			handleErr,
		),
		cacheprefetch.WeakCache(
			cacheKey,
			typeStr,
			nil,
			cache.SloshTrue[*M],
			nil,
		),
	)
}

func infoGetFetchFunc[M any](
	url, method, cacheKey, typeStr string,
	minExpiresIn time.Duration,
	auth *RefreshTokenAndApp,
	handleErr func(error) (ok bool, rep *M, expires time.Time),
) fetch.CachingFetch[*M] {
	return func(x cache.Context) (
		rep *M,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		rep, expires, err = fetch.FetchWithRetries(
			x,
			func(x cache.Context) (rep *M, expires time.Time, err error) {
				var repVal M
				repVal, expires, err = getModel[M](x, url, method, auth, nil)
				if err != nil {
					var ok bool
					ok, rep, expires = handleErr(err)
					if !ok {
						return nil, expires, err
					}
				} else {
					rep = &repVal
				}
				return rep, expires, nil
			},
			ESI_NUM_RETRIES,
			esiShouldRetry,
		)
		if err != nil {
			return nil, expires, nil, err
		}
		expires = fetch.CalcExpiresIn(expires, minExpiresIn)
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.DualSetOne[*M](
				cacheKey,
				typeStr,
				rep,
				expires,
			),
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
			rep = ForbiddenStructure
		} else if statusErr.Code == http.StatusForbidden {
			ok = true
			rep = ForbiddenStructure
		}
	}
	return ok, rep, expires
}
