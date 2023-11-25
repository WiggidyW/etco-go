package service

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
)

func authorizedGet[T any](
	x cache.Context,
	refreshToken string,
	getAuthorized func(cache.Context, string) (bool, time.Time, error),
	getRep func(cache.Context) (T, time.Time, error),
) (
	authorized bool,
	empty T,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	chnAuthorized := expirable.NewChanResult[bool](x.Ctx(), 1, 0)
	go expirable.P2Transceive(
		chnAuthorized,
		x, refreshToken,
		getAuthorized,
	)

	var rep T
	rep, _, err = getRep(x)
	if err != nil {
		return false, empty, err
	}

	authorized, _, err = chnAuthorized.RecvExp()
	if !authorized || err != nil {
		return false, empty, err
	}

	return true, rep, nil
}

func authorizedGetP1[T any, P1 any](
	x cache.Context,
	refreshToken string,
	getAuthorized func(cache.Context, string) (bool, time.Time, error),
	getRep func(
		cache.Context,
		P1,
	) (T, time.Time, error),
	p1 P1,
) (
	authorized bool,
	empty T,
	err error,
) {
	return authorizedGet[T](
		x, refreshToken, getAuthorized,
		func(cache.Context) (T, time.Time, error) {
			return getRep(x, p1)
		},
	)
}

func authorizedGetP2[T any, P1 any, P2 any](
	x cache.Context,
	refreshToken string,
	getAuthorized func(cache.Context, string) (bool, time.Time, error),
	getRep func(
		cache.Context,
		P1,
		P2,
	) (T, time.Time, error),
	p1 P1,
	p2 P2,
) (
	authorized bool,
	empty T,
	err error,
) {
	return authorizedGet[T](
		x, refreshToken, getAuthorized,
		func(cache.Context) (T, time.Time, error) {
			return getRep(x, p1, p2)
		},
	)
}

// func authorizedGetP3[T any, P1 any, P2 any, P3 any](
// 	x cache.Context,
// 	refreshToken string,
// 	getAuthorized func(cache.Context, string) (bool, time.Time, error),
// 	getRep func(
// 		cache.Context,
// 		P1,
// 		P2,
// 		P3,
// 	) (T, time.Time, error),
// 	p1 P1,
// 	p2 P2,
// 	p3 P3,
// ) (
// 	authorized bool,
// 	empty T,
// 	err error,
// ) {
// 	return authorizedGet[T](
// 		x, refreshToken, getAuthorized,
// 		func(cache.Context) (T, time.Time, error) {
// 			return getRep(x, p1, p2, p3)
// 		},
// 	)
// }
