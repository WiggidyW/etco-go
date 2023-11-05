package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
)

const (
	USERDATA_BUF_CAP          int           = 0
	USERDATA_LOCK_TTL         time.Duration = 1 * time.Minute
	USERDATA_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	USERDATA_EXPIRES_IN       time.Duration = 24 * time.Hour

	USER_B_CODES_BUF_CAP          int           = 0
	USER_B_CODES_LOCK_TTL         time.Duration = 1 * time.Minute
	USER_B_CODES_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute

	USER_S_CODES_BUF_CAP          int           = 0
	USER_S_CODES_LOCK_TTL         time.Duration = 1 * time.Minute
	USER_S_CODES_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute

	USER_C_PURCHASE_BUF_CAP          int           = 0
	USER_C_PURCHASE_LOCK_TTL         time.Duration = 1 * time.Minute
	USER_C_PURCHASE_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute

	USER_M_PURCHASE_BUF_CAP          int           = 0
	USER_M_PURCHASE_LOCK_TTL         time.Duration = 1 * time.Minute
	USER_M_PURCHASE_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

func init() {
	keys.TypeStrUserData = localcache.RegisterType[UserData](USERDATA_BUF_CAP)
	keys.TypeStrUserBuybackAppraisalCodes = localcache.RegisterType[[]string](USER_B_CODES_BUF_CAP)
	keys.TypeStrUserShopAppraisalCodes = localcache.RegisterType[[]string](USER_S_CODES_BUF_CAP)
	keys.TypeStrUserCancelledPurchase = localcache.RegisterType[*time.Time](USER_C_PURCHASE_BUF_CAP)
	keys.TypeStrUserMadePurchase = localcache.RegisterType[*time.Time](USER_M_PURCHASE_BUF_CAP)
}

type UserData struct {
	BuybackAppraisals []string   `firestore:"buyback_appraisals"`
	ShopAppraisals    []string   `firestore:"shop_appraisals"`
	CancelledPurchase *time.Time `firestore:"cancelled_purchase"`
	MadePurchase      *time.Time `firestore:"made_purchase"`
}

func GetUserData(
	ctx context.Context,
	characterId int32,
) (
	rep UserData,
	expires *time.Time,
	err error,
) {
	return userDataGet(
		ctx,
		characterId,
		USERDATA_LOCK_TTL,
		USERDATA_LOCK_MAX_BACKOFF,
		USERDATA_EXPIRES_IN,
	)
}

func GetUserBuybackAppraisalCodes(
	ctx context.Context,
	characterId int32,
) (
	rep []string,
	expires *time.Time,
	err error,
) {
	return userDataFieldGet(
		ctx,
		characterId,
		keys.TypeStrUserBuybackAppraisalCodes,
		keys.CacheKeyUserBuybackAppraisalCodes(characterId),
		USER_B_CODES_LOCK_TTL,
		USER_B_CODES_LOCK_MAX_BACKOFF,
		func(userData UserData) *[]string {
			return &userData.BuybackAppraisals
		},
	)
}

func GetUserShopAppraisalCodes(
	ctx context.Context,
	characterId int32,
) (
	rep []string,
	expires *time.Time,
	err error,
) {
	return userDataFieldGet(
		ctx,
		characterId,
		keys.TypeStrUserShopAppraisalCodes,
		keys.CacheKeyUserShopAppraisalCodes(characterId),
		USER_S_CODES_LOCK_TTL,
		USER_S_CODES_LOCK_MAX_BACKOFF,
		func(userData UserData) *[]string {
			return &userData.ShopAppraisals
		},
	)
}

func GetUserMadePurchase(
	ctx context.Context,
	characterId int32,
) (
	rep *time.Time,
	expires *time.Time,
	err error,
) {
	return userDataFieldGet(
		ctx,
		characterId,
		keys.TypeStrUserMadePurchase,
		keys.CacheKeyUserMadePurchase(characterId),
		USER_M_PURCHASE_LOCK_TTL,
		USER_M_PURCHASE_LOCK_MAX_BACKOFF,
		func(userData UserData) **time.Time {
			return &userData.MadePurchase
		},
	)
}

func GetUserCancelledPurchase(
	ctx context.Context,
	characterId int32,
) (
	rep *time.Time,
	expires *time.Time,
	err error,
) {
	return userDataFieldGet(
		ctx,
		characterId,
		keys.TypeStrUserCancelledPurchase,
		keys.CacheKeyUserCancelledPurchase(characterId),
		USER_C_PURCHASE_LOCK_TTL,
		USER_C_PURCHASE_LOCK_MAX_BACKOFF,
		func(userData UserData) **time.Time {
			return &userData.CancelledPurchase
		},
	)
}
