package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	USERDATA_EXPIRES_IN     time.Duration = 24 * time.Hour
	USERDATA_BUF_CAP        int           = 0
	USER_B_CODES_BUF_CAP    int           = 0
	USER_S_CODES_BUF_CAP    int           = 0
	USER_C_PURCHASE_BUF_CAP int           = 0
	USER_M_PURCHASE_BUF_CAP int           = 0
)

func init() {
	keys.TypeStrNSUserData = cache.RegisterType[UserData]("userdata", USERDATA_BUF_CAP)
	keys.TypeStrUserBuybackAppraisalCodes = cache.RegisterType[[]string]("userbuybackappraisalcodes", USER_B_CODES_BUF_CAP)
	keys.TypeStrUserShopAppraisalCodes = cache.RegisterType[[]string]("usershopappraisalcodes", USER_S_CODES_BUF_CAP)
	keys.TypeStrUserCancelledPurchase = cache.RegisterType[*time.Time]("usercancelledpurchase", USER_C_PURCHASE_BUF_CAP)
	keys.TypeStrUserMadePurchase = cache.RegisterType[*time.Time]("usermadepurchase", USER_M_PURCHASE_BUF_CAP)
}

type UserData struct {
	BuybackAppraisals []string   `firestore:"buyback_appraisals"`
	ShopAppraisals    []string   `firestore:"shop_appraisals"`
	CancelledPurchase *time.Time `firestore:"cancelled_purchase"`
	MadePurchase      *time.Time `firestore:"made_purchase"`
}

func GetUserBuybackAppraisalCodes(
	x cache.Context,
	characterId int32,
) (
	rep []string,
	expires time.Time,
	err error,
) {
	return userDataFieldGet(
		x,
		characterId,
		udf_B_APPRAISAL_CODES,
		func(userData UserData) *[]string {
			return &userData.BuybackAppraisals
		},
	)
}

func GetUserShopAppraisalCodes(
	x cache.Context,
	characterId int32,
) (
	rep []string,
	expires time.Time,
	err error,
) {
	return userDataFieldGet(
		x,
		characterId,
		udf_S_APPRAISAL_CODES,
		func(userData UserData) *[]string {
			return &userData.ShopAppraisals
		},
	)
}

func GetUserCancelledPurchase(
	x cache.Context,
	characterId int32,
) (
	rep *time.Time,
	expires time.Time,
	err error,
) {
	return userDataFieldGet(
		x,
		characterId,
		udf_C_PURCHASE,
		func(userData UserData) **time.Time {
			return &userData.CancelledPurchase
		},
	)
}

func GetUserMadePurchase(
	x cache.Context,
	characterId int32,
) (
	rep *time.Time,
	expires time.Time,
	err error,
) {
	return userDataFieldGet(
		x,
		characterId,
		udf_M_PURCHASE,
		func(userData UserData) **time.Time {
			return &userData.MadePurchase
		},
	)
}
