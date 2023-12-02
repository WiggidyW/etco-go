package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"
	"github.com/WiggidyW/etco-go/util"
)

const (
	PREV_CONTRACTS_EXPIRES_IN time.Duration = 24 * time.Hour
	PREV_CONTRACTS_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrPrevContracts = cache.RegisterType[PreviousContracts]("prevcontracts", PREV_CONTRACTS_BUF_CAP)
}

type PreviousContracts = implrdb.PreviousContracts

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
			rep, err = readPrevContracts(x.Ctx())
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
			err = setPrevContracts(x.Ctx(), rep.Buyback, rep.Shop, rep.Haul)
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

type NewContracts[BV any, SV any, HV any] struct {
	Buyback map[string]BV
	Shop    map[string]SV
	Haul    map[string]HV
}

// Using the parameters, which are the current contracts, we will return
// all contracts that have not been seen before.
//
// Possible Side Effect - Sets PreviousContracts (if there is a diff)
//
// ^^^ will happen unless the rep is empty. Even if rep is empty, sometimes
// it will still happen (when current contracts has no new, but is fewer)
func GetNewContracts[BV any, SV any, HV any](
	x cache.Context,
	buybackContracts map[string]BV,
	shopContracts map[string]SV,
	haulContracts map[string]HV,
) (
	rep NewContracts[BV, SV, HV],
	expires time.Time,
	err error,
) {
	var prevContracts PreviousContracts
	prevContracts, expires, err = GetPrevContracts(x)
	if err != nil {
		return rep, expires, err
	}
	newBuybackContracts := util.KeysNotIn(
		buybackContracts,
		util.SliceToSet(prevContracts.Buyback),
	)
	newShopContracts := util.KeysNotIn(
		shopContracts,
		util.SliceToSet(prevContracts.Shop),
	)
	newHaulContracts := util.KeysNotIn(
		haulContracts,
		util.SliceToSet(prevContracts.Haul),
	)

	// set prev contracts if:
	// - there are new contracts (any new code doesn't exist in prev)
	// - there is a length difference (new is a subset of prev with nothing new)
	//
	// If there are no new codes, and the length is the same, then
	// previous contracts == current contracts.
	if len(newBuybackContracts) != 0 ||
		len(newShopContracts) != 0 ||
		len(newHaulContracts) != 0 ||
		len(buybackContracts) != len(prevContracts.Buyback) ||
		len(shopContracts) != len(prevContracts.Shop) ||
		len(haulContracts) != len(prevContracts.Haul) {
		go logger.MaybeErr(SetPrevContracts(
			x,
			PreviousContracts{
				Buyback: util.KeysToSlice(buybackContracts),
				Shop:    util.KeysToSlice(shopContracts),
				Haul:    util.KeysToSlice(haulContracts),
			},
		))
	}

	rep = NewContracts[BV, SV, HV]{
		Buyback: newBuybackContracts,
		Shop:    newShopContracts,
		Haul:    newHaulContracts,
	}
	return rep, expires, nil
}
