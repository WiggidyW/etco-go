package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

type userDataField uint8

const (
	udf_B_APPRAISAL_CODES userDataField = iota
	udf_S_APPRAISAL_CODES
	udf_M_PURCHASE
	udf_C_PURCHASE
)

type userDataKeys struct {
	NSCacheKey              string
	NSTypeStr               string
	BAppraisalCodesCacheKey string
	BAppraisalCodesTypeStr  string
	SAppraisalCodesCacheKey string
	SAppraisalCodesTypeStr  string
	CPurchaseCacheKey       string
	CPurchaseTypeStr        string
	MPurchaseCacheKey       string
	MPurchaseTypeStr        string
}

func newUserDataKeysAndGetKeys(
	characterId int32,
	getKind userDataField,
) (
	k userDataKeys,
	getCacheKey string,
	getTypeStr string,
) {
	k = userDataKeys{
		NSCacheKey:              keys.CacheKeyNSUserData(characterId),
		NSTypeStr:               keys.TypeStrNSUserData,
		BAppraisalCodesCacheKey: keys.CacheKeyUserBuybackAppraisalCodes(characterId),
		BAppraisalCodesTypeStr:  keys.TypeStrUserBuybackAppraisalCodes,
		SAppraisalCodesCacheKey: keys.CacheKeyUserShopAppraisalCodes(characterId),
		SAppraisalCodesTypeStr:  keys.TypeStrUserShopAppraisalCodes,
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
	return fetch.HandleFetch[T](
		x,
		&prefetch.Params[T]{
			CacheParams: &prefetch.CacheParams[T]{
				Get: prefetch.ServerCacheGet[T](
					cacheKey, typeStr,
					false,
					nil,
				),
				Namespace: prefetch.CacheNamespace(
					k.NSCacheKey,
					k.NSTypeStr,
					false,
				),
				Lock: prefetch.CacheOrderedLocksNoFamily(
					nil,
					prefetch.ServerCacheLock(
						k.BAppraisalCodesCacheKey,
						k.BAppraisalCodesTypeStr,
					),
					prefetch.ServerCacheLock(
						k.SAppraisalCodesCacheKey,
						k.SAppraisalCodesTypeStr,
					),
					prefetch.ServerCacheLock(
						k.CPurchaseCacheKey,
						k.CPurchaseTypeStr,
					),
					prefetch.ServerCacheLock(
						k.MPurchaseCacheKey,
						k.MPurchaseTypeStr,
					),
				),
			},
		},
		userDataFieldGetFetchFunc[T](
			characterId,
			k,
			getField,
		),
		nil,
	)
}

func userDataFieldGetFetchFunc[T any](
	characterId int32,
	keys userDataKeys,
	getField func(UserData) T,
) fetch.Fetch[T] {
	return func(x cache.Context) (
		rep T,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var userData UserData
		userData, err = client.readUserData(x.Ctx(), characterId)
		if err != nil {
			return rep, expires, nil, err
		}
		expires = time.Now().Add(USERDATA_EXPIRES_IN)
		rep = getField(userData)
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Namespace: postfetch.CacheNamespace(
					keys.NSCacheKey,
					keys.NSTypeStr,
					expires,
				),
				Set: []postfetch.CacheActionSet{
					postfetch.ServerCacheSet(
						keys.BAppraisalCodesCacheKey,
						keys.BAppraisalCodesTypeStr,
						&userData.BuybackAppraisals,
						expires,
					),
					postfetch.ServerCacheSet(
						keys.SAppraisalCodesCacheKey,
						keys.SAppraisalCodesTypeStr,
						&userData.ShopAppraisals,
						expires,
					),
					postfetch.ServerCacheSet(
						keys.CPurchaseCacheKey,
						keys.CPurchaseTypeStr,
						&userData.CancelledPurchase,
						expires,
					),
					postfetch.ServerCacheSet(
						keys.MPurchaseCacheKey,
						keys.MPurchaseTypeStr,
						&userData.MadePurchase,
						expires,
					),
				},
			},
		}
		return rep, expires, postFetch, nil
	}
}
