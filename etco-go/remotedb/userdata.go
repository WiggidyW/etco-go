package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"

	// "github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	USERDATA_EXPIRES_IN     time.Duration = 24 * time.Hour
	USERDATA_BUF_CAP        int           = 0
	USER_B_CODES_BUF_CAP    int           = 0
	USER_S_CODES_BUF_CAP    int           = 0
	USER_H_CODES_BUF_CAP    int           = 0
	USER_C_PURCHASE_BUF_CAP int           = 0
	USER_M_PURCHASE_BUF_CAP int           = 0
)

func init() {
	keys.TypeStrNSUserData = cache.RegisterType[UserData]("userdata", USERDATA_BUF_CAP)
	keys.TypeStrUserBuybackAppraisalCodes = cache.RegisterType[[]string]("userbuybackappraisalcodes", USER_B_CODES_BUF_CAP)
	keys.TypeStrUserShopAppraisalCodes = cache.RegisterType[[]string]("usershopappraisalcodes", USER_S_CODES_BUF_CAP)
	keys.TypeStrUserHaulAppraisalCodes = cache.RegisterType[[]string]("userhaulappraisalcodes", USER_H_CODES_BUF_CAP)
	// keys.TypeStrUserCancelledPurchase = cache.RegisterType[*time.Time]("usercancelledpurchase", USER_C_PURCHASE_BUF_CAP)
	// keys.TypeStrUserMadePurchase = cache.RegisterType[*time.Time]("usermadepurchase", USER_M_PURCHASE_BUF_CAP)
	keys.TypeStrUserCancelledPurchase = cache.RegisterType[time.Time]("usercancelledpurchase", USER_C_PURCHASE_BUF_CAP)
	keys.TypeStrUserMadePurchase = cache.RegisterType[time.Time]("usermadepurchase", USER_M_PURCHASE_BUF_CAP)
}

type UserData = implrdb.UserData

func udGetBuybackAppraisals(ud UserData) []string   { return ud.BuybackAppraisals }
func udGetShopAppraisals(ud UserData) []string      { return ud.ShopAppraisals }
func udGetHaulAppraisals(ud UserData) []string      { return ud.HaulAppraisals }
func udGetCancelledPurchase(ud UserData) *time.Time { return ud.CancelledPurchase }
func udGetMadePurchase(ud UserData) *time.Time      { return ud.MadePurchase }

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
		udGetBuybackAppraisals,
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
		udGetShopAppraisals,
	)
}

func GetUserHaulAppraisalCodes(
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
		udf_H_APPRAISAL_CODES,
		udGetHaulAppraisals,
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
		udGetCancelledPurchase,
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
		udGetMadePurchase,
	)
}
