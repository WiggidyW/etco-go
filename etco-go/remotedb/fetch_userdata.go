package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func userDataGet(
	ctx context.Context,
	characterId int32,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
) (
	rep UserData,
	expires *time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyUserData(characterId)

	var repPtr *UserData
	repPtr, expires, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[UserData]{
			CacheParams: &prefetch.CacheParams[UserData]{
				Get: prefetch.ServerCacheGet[UserData](
					keys.TypeStrUserData, cacheKey,
					lockTTL, lockMaxBackoff,
					nil,
				),
			},
		},
		userDataGetFetchFunc(characterId, cacheKey, expiresIn),
	)
	if err != nil {
		return rep, nil, err
	} else if repPtr != nil {
		rep = *repPtr
	}
	return rep, expires, nil
}

func userDataGetFetchFunc(
	characterId int32,
	cacheKey string,
	expiresIn time.Duration,
) fetch.Fetch[UserData] {
	return func(ctx context.Context) (
		rep *UserData,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		rep, err = client.readUserData(ctx, characterId)
		if err != nil {
			return nil, nil, nil, err
		}
		expires = fetch.ExpiresIn(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.ServerCacheSet(
				keys.TypeStrUserData, cacheKey,
			),
		}
		return rep, expires, postFetch, nil
	}
}

func userDataFieldGet[T any](
	ctx context.Context,
	characterId int32,
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
	getField func(UserData) *T,
) (
	rep T,
	expires *time.Time,
	err error,
) {
	var repPtr *T
	repPtr, expires, err = fetch.HandleFetch[T](
		ctx,
		&prefetch.Params[T]{
			CacheParams: &prefetch.CacheParams[T]{
				Get: prefetch.ServerCacheGet[T](
					typeStr, cacheKey,
					lockTTL, lockMaxBackoff,
					nil,
				),
			},
		},
		userDataFieldGetFetchFunc[T](
			characterId,
			typeStr, cacheKey,
			getField,
		),
	)
	if err != nil {
		return rep, nil, err
	} else if repPtr != nil {
		rep = *repPtr
	}
	return rep, expires, nil
}

func userDataFieldGetFetchFunc[T any](
	characterId int32,
	typeStr, cacheKey string,
	getField func(UserData) *T,
) fetch.Fetch[T] {
	return func(ctx context.Context) (
		rep *T,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var userData UserData
		userData, expires, err = GetUserData(ctx, characterId)
		if err != nil {
			return nil, nil, nil, err
		}
		rep = getField(userData)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.ServerCacheSet(typeStr, cacheKey),
		}
		return rep, expires, postFetch, nil
	}
}
