package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
)

type userDataField uint8

const (
	udf_B_APPRAISAL_CODES userDataField = iota
	udf_S_APPRAISAL_CODES
	udf_H_APPRAISAL_CODES
	udf_M_PURCHASE
	udf_C_PURCHASE
)

type userDataKeys struct {
	NSCacheKey              keys.Key
	NSTypeStr               keys.Key
	BAppraisalCodesCacheKey keys.Key
	BAppraisalCodesTypeStr  keys.Key
	SAppraisalCodesCacheKey keys.Key
	SAppraisalCodesTypeStr  keys.Key
	HAppraisalCodesCacheKey keys.Key
	HAppraisalCodesTypeStr  keys.Key
	CPurchaseCacheKey       keys.Key
	CPurchaseTypeStr        keys.Key
	MPurchaseCacheKey       keys.Key
	MPurchaseTypeStr        keys.Key
}

func newUserDataKeysAndGetKeys(
	characterId int32,
	getKind userDataField,
) (
	k userDataKeys,
	getCacheKey keys.Key,
	getTypeStr keys.Key,
) {
	k = userDataKeys{
		NSCacheKey:              keys.CacheKeyNSUserData(characterId),
		NSTypeStr:               keys.TypeStrNSUserData,
		BAppraisalCodesCacheKey: keys.CacheKeyUserBuybackAppraisalCodes(characterId),
		BAppraisalCodesTypeStr:  keys.TypeStrUserBuybackAppraisalCodes,
		SAppraisalCodesCacheKey: keys.CacheKeyUserShopAppraisalCodes(characterId),
		SAppraisalCodesTypeStr:  keys.TypeStrUserShopAppraisalCodes,
		HAppraisalCodesCacheKey: keys.CacheKeyUserHaulAppraisalCodes(characterId),
		HAppraisalCodesTypeStr:  keys.TypeStrUserHaulAppraisalCodes,
		CPurchaseCacheKey:       keys.CacheKeyUserCancelledPurchase(characterId),
		CPurchaseTypeStr:        keys.TypeStrUserCancelledPurchase,
		MPurchaseCacheKey:       keys.CacheKeyUserMadePurchase(characterId),
		MPurchaseTypeStr:        keys.TypeStrUserMadePurchase,
	}
	switch getKind {
	case udf_B_APPRAISAL_CODES:
		getCacheKey = k.BAppraisalCodesCacheKey
		getTypeStr = k.BAppraisalCodesTypeStr
	case udf_S_APPRAISAL_CODES:
		getCacheKey = k.SAppraisalCodesCacheKey
		getTypeStr = k.SAppraisalCodesTypeStr
	case udf_H_APPRAISAL_CODES:
		getCacheKey = k.HAppraisalCodesCacheKey
		getTypeStr = k.HAppraisalCodesTypeStr
	case udf_C_PURCHASE:
		getCacheKey = k.CPurchaseCacheKey
		getTypeStr = k.CPurchaseTypeStr
	case udf_M_PURCHASE:
		getCacheKey = k.MPurchaseCacheKey
		getTypeStr = k.MPurchaseTypeStr
	default:
		panic("invalid getKind")
	}
	return k, getCacheKey, getTypeStr
}

func userDataFieldGet[T any](
	x cache.Context,
	characterId int32,
	getKind userDataField,
	getField func(UserData) T,
) (
	rep T,
	expires time.Time,
	err error,
) {
	k, cacheKey, typeStr := newUserDataKeysAndGetKeys(
		characterId,
		getKind,
	)
	return fetch.FetchWithCache[T](
		x,
		userDataFieldGetFetchFunc[T](
			characterId,
			k,
			getField,
		),
		cacheprefetch.StrongMultiCacheKnownKeys[T](
			cacheKey,
			typeStr,
			k.NSCacheKey,
			k.NSTypeStr,
			nil,
			[]cacheprefetch.ActionOrderedLocks{{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.ServerLock(
						k.BAppraisalCodesCacheKey,
						k.BAppraisalCodesTypeStr,
					),
					cacheprefetch.ServerLock(
						k.SAppraisalCodesCacheKey,
						k.SAppraisalCodesTypeStr,
					),
					cacheprefetch.ServerLock(
						k.HAppraisalCodesCacheKey,
						k.HAppraisalCodesTypeStr,
					),
					cacheprefetch.ServerLock(
						k.CPurchaseCacheKey,
						k.CPurchaseTypeStr,
					),
					cacheprefetch.ServerLock(
						k.MPurchaseCacheKey,
						k.MPurchaseTypeStr,
					),
				},
				Child: nil,
			}},
		),
	)
}

func userDataFieldGetFetchFunc[T any](
	characterId int32,
	keys userDataKeys,
	getField func(UserData) T,
) fetch.CachingFetch[T] {
	return func(x cache.Context) (
		rep T,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		var userData UserData
		userData, err = readUserData(x.Ctx(), characterId)
		if err != nil {
			return rep, expires, nil, err
		}
		expires = time.Now().Add(USERDATA_EXPIRES_IN)
		rep = getField(userData)
		postFetch = &cachepostfetch.Params{
			Namespace: cachepostfetch.Namespace(
				keys.NSCacheKey,
				keys.NSTypeStr,
				expires,
			),
			Set: []cachepostfetch.ActionSet{
				cachepostfetch.ServerSet[[]string](
					keys.BAppraisalCodesCacheKey,
					keys.BAppraisalCodesTypeStr,
					userData.BuybackAppraisals,
					expires,
				),
				cachepostfetch.ServerSet[[]string](
					keys.SAppraisalCodesCacheKey,
					keys.SAppraisalCodesTypeStr,
					userData.ShopAppraisals,
					expires,
				),
				cachepostfetch.ServerSet[[]string](
					keys.HAppraisalCodesCacheKey,
					keys.HAppraisalCodesTypeStr,
					userData.HaulAppraisals,
					expires,
				),
				cachepostfetch.ServerSet[*time.Time](
					keys.CPurchaseCacheKey,
					keys.CPurchaseTypeStr,
					userData.CancelledPurchase,
					expires,
				),
				cachepostfetch.ServerSet[*time.Time](
					keys.MPurchaseCacheKey,
					keys.MPurchaseTypeStr,
					userData.MadePurchase,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}
