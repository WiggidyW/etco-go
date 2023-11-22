package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/util"
)

const (
	PREV_CONTRACTS_EXPIRES_IN time.Duration = 24 * time.Hour
	PREV_CONTRACTS_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrPrevContracts = cache.RegisterType[PreviousContracts]("prevcontracts", PREV_CONTRACTS_BUF_CAP)
}

type PreviousContracts struct {
	Buyback []string `firestore:"buyback_codes"`
	Shop    []string `firestore:"shop_codes"`
}

func GetPrevContracts(x cache.Context) (
	rep PreviousContracts,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache(
		x,
		func(x cache.Context) (
			rep PreviousContracts,
			expires time.Time,
			postFetch *cachepostfetch.Params,
			err error,
		) {
			rep, err = client.readPrevContracts(x.Ctx())
			if err != nil {
				return rep, expires, nil, err
			}
			expires = time.Now().Add(PREV_CONTRACTS_EXPIRES_IN)
			postFetch = &cachepostfetch.Params{
				Set: cachepostfetch.ServerSetOne[PreviousContracts](
					keys.CacheKeyPrevContracts,
					keys.TypeStrPrevContracts,
					rep,
					expires,
				),
			}
			return rep, expires, postFetch, nil
		},
		cacheprefetch.StrongCache[PreviousContracts](
			keys.CacheKeyPrevContracts,
			keys.TypeStrPrevContracts,
			nil,
			nil,
		),
	)
}

func SetPrevContracts(
	x cache.Context,
	rep PreviousContracts,
) (
	err error,
) {
	_, _, err = fetch.FetchWithCache(
		x,
		func(x cache.Context) (
			_ struct{},
			_ time.Time,
			postFetch *cachepostfetch.Params,
			err error,
		) {
			err = client.setPrevContracts(x.Ctx(), rep.Buyback, rep.Shop)
			if err != nil {
				return struct{}{}, time.Time{}, nil, err
			}
			expires := time.Now().Add(PREV_CONTRACTS_EXPIRES_IN)
			postFetch = &cachepostfetch.Params{
				Set: cachepostfetch.ServerSetOne[PreviousContracts](
					keys.CacheKeyPrevContracts,
					keys.TypeStrPrevContracts,
					rep,
					expires,
				),
			}
			return struct{}{}, expires, postFetch, nil
		},
		cacheprefetch.AntiCache[struct{}]([]cacheprefetch.ActionOrderedLocks{{
			Locks: []cacheprefetch.ActionLock{
				cacheprefetch.ServerLock(
					keys.CacheKeyPrevContracts,
					keys.TypeStrPrevContracts,
				),
			},
			Child: nil,
		}}),
	)
	return err
}

type NewContracts = PreviousContracts

// Side Effect - Sets PreviousContracts
func GetNewContracts[BV any, SV any](
	x cache.Context,
	buybackContractsMap map[string]BV,
	shopContractsMap map[string]SV,
) (
	rep NewContracts,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache(
		x,
		func(x cache.Context) (
			rep NewContracts,
			expires time.Time,
			_ *cachepostfetch.Params,
			err error,
		) {
			var prevContracts PreviousContracts
			prevContracts, expires, err = GetPrevContracts(x)
			if err != nil {
				return rep, expires, nil, err
			}
			buybackContracts := util.KeysToSlice(buybackContractsMap)
			shopContracts := util.KeysToSlice(shopContractsMap)
			newBuybackContracts := util.KeysNotIn(
				buybackContracts,
				util.SliceToSet(prevContracts.Buyback),
			)
			newShopContracts := util.KeysNotIn(
				shopContracts,
				util.SliceToSet(prevContracts.Shop),
			)
			go logger.MaybeErr(SetPrevContracts(
				x,
				PreviousContracts{
					Buyback: buybackContracts,
					Shop:    shopContracts,
				},
			))
			return NewContracts{
				Buyback: newBuybackContracts,
				Shop:    newShopContracts,
			}, expires, nil, nil
		},
		cacheprefetch.AntiCache[NewContracts](
			[]cacheprefetch.ActionOrderedLocks{{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.ServerLock(
						keys.CacheKeyPrevContracts,
						keys.TypeStrPrevContracts,
					),
				},
				Child: nil,
			},
			}),
	)
}
