package builder

import (
	"fmt"
	"os"
	"strconv"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go-builder/builderenv"
)

// func transceiveWriteConstants(
// 	filePath string,
// 	updaterBucketData b.UpdaterData,
// 	bootstrapAdminId int32,
// 	corporationId int32,
// 	corporationWebRefreshToken string,
// 	structureInfoWebRefreshToken string,
// 	remotedbProjectId string,
// 	bucketCredsJson string,
// 	remotedbCredsJson string,
// 	chnSend chanresult.ChanSendResult[struct{}],
// ) error {
// 	err := writeConstants(
// 		filePath,
// 		updaterBucketData,
// 		bootstrapAdminId,
// 		corporationId,
// 		corporationWebRefreshToken,
// 		structureInfoWebRefreshToken,
// 		remotedbProjectId,
// 		bucketCredsJson,
// 		remotedbCredsJson,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(struct{}{})
// 	}
// }

func writeConstants(
	filePath string,
	constantsBucketData b.ConstantsData,
	updaterBucketData b.UpdaterData,
) error {
	// If any are missing from bucket data, set them to ENV values.
	if constantsBucketData.PURCHASE_MAX_ACTIVE == nil {
		constantsBucketData.PURCHASE_MAX_ACTIVE = &builderenv.
			PURCHASE_MAX_ACTIVE
	}
	if constantsBucketData.MAKE_PURCHASE_COOLDOWN == nil {
		constantsBucketData.MAKE_PURCHASE_COOLDOWN = &builderenv.
			MAKE_PURCHASE_COOLDOWN
	}
	if constantsBucketData.CANCEL_PURCHASE_COOLDOWN == nil {
		constantsBucketData.CANCEL_PURCHASE_COOLDOWN = &builderenv.
			CANCEL_PURCHASE_COOLDOWN
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	constantsFileString := fmt.Sprintf(
		`package buildconstants

		import "time"

		const (
			CAPACITY_SDE_CATEGORIES      int = %d
			CAPACITY_SDE_GROUPS          int = %d
			CAPACITY_SDE_MARKET_GROUPS   int = %d
			CAPACITY_SDE_NAME_TO_TYPE_ID int = %d
			CAPACITY_SDE_REGIONS         int = %d
			CAPACITY_SDE_SYSTEMS         int = %d
			CAPACITY_SDE_STATIONS        int = %d
			CAPACITY_SDE_TYPE_DATA_MAP   int = %d
			CAPACITY_SDE_TYPE_VOLUMES    int = %d

			CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER int = %d
			CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER  int = %d
			CAPACITY_WEB_BUYBACK_SYSTEMS                  int = %d
			CAPACITY_WEB_SHOP_LOCATIONS                   int = %d
			CAPACITY_WEB_MARKETS                          int = %d

			CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS int = %d
			CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS  int = %d
			CAPACITY_CORE_BUYBACK_SYSTEMS          int = %d
			CAPACITY_CORE_SHOP_LOCATIONS           int = %d
			CAPACITY_CORE_BANNED_FLAG_SETS         int = %d
			CAPACITY_CORE_PRICINGS                 int = %d
			CAPACITY_CORE_MARKETS                  int = %d

			VERSION_BUYBACK string = "%s"
			VERSION_SHOP    string = "%s"

			CCACHE_MAX_BYTES int    = %d
			SCACHE_ADDRESS   string = "%s"

			BOOTSTRAP_ADMIN_ID               int32  = %d
			CORPORATION_ID                   int32  = %d
			CORPORATION_WEB_REFRESH_TOKEN    string = "%s"
			STRUCTURE_INFO_WEB_REFRESH_TOKEN string = "%s"

			REMOTEDB_PROJECT_ID string = "%s"
			REMOTEDB_CREDS_JSON string = %s
			BUCKET_CREDS_JSON   string = %s

			PURCHASE_MAX_ACTIVE      int = %d
			MAKE_PURCHASE_COOLDOWN   time.Duration = %d
			CANCEL_PURCHASE_COOLDOWN time.Duration = %d

			ESI_USER_AGENT                   string = "%s"
			ESI_MARKETS_CLIENT_ID            string = "%s"
			ESI_MARKETS_CLIENT_SECRET        string = "%s"
			ESI_CORP_CLIENT_ID               string = "%s"
			ESI_CORP_CLIENT_SECRET           string = "%s"
			ESI_STRUCTURE_INFO_CLIENT_ID     string = "%s"
			ESI_STRUCTURE_INFO_CLIENT_SECRET string = "%s"
			ESI_AUTH_CLIENT_ID               string = "%s"
			ESI_AUTH_CLIENT_SECRET           string = "%s"
		)
		`,
		updaterBucketData.CAPACITY_SDE_CATEGORIES,
		updaterBucketData.CAPACITY_SDE_GROUPS,
		updaterBucketData.CAPACITY_SDE_MARKET_GROUPS,
		updaterBucketData.CAPACITY_SDE_NAME_TO_TYPE_ID,
		updaterBucketData.CAPACITY_SDE_REGIONS,
		updaterBucketData.CAPACITY_SDE_SYSTEMS,
		updaterBucketData.CAPACITY_SDE_STATIONS,
		updaterBucketData.CAPACITY_SDE_TYPE_DATA_MAP,
		updaterBucketData.CAPACITY_SDE_TYPE_VOLUMES,

		updaterBucketData.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
		updaterBucketData.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
		updaterBucketData.CAPACITY_WEB_BUYBACK_SYSTEMS,
		updaterBucketData.CAPACITY_WEB_SHOP_LOCATIONS,
		updaterBucketData.CAPACITY_WEB_MARKETS,

		updaterBucketData.CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS,
		updaterBucketData.CAPACITY_CORE_SHOP_LOCATIONS,
		updaterBucketData.CAPACITY_CORE_BANNED_FLAG_SETS,
		updaterBucketData.CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS,
		updaterBucketData.CAPACITY_CORE_BUYBACK_SYSTEMS,
		updaterBucketData.CAPACITY_CORE_PRICINGS,
		updaterBucketData.CAPACITY_CORE_MARKETS,

		updaterBucketData.VERSION_BUYBACK,
		updaterBucketData.VERSION_SHOP,

		builderenv.CCACHE_MAX_BYTES,
		builderenv.SCACHE_ADDRESS,

		builderenv.BOOTSTRAP_ADMIN_ID,
		builderenv.CORPORATION_ID,
		builderenv.CORPORATION_WEB_REFRESH_TOKEN,
		builderenv.STRUCTURE_INFO_WEB_REFRESH_TOKEN,

		builderenv.REMOTEDB_PROJECT_ID,
		strconv.Quote(builderenv.REMOTEDB_CREDS_JSON),
		strconv.Quote(builderenv.BUCKET_CREDS_JSON),

		*constantsBucketData.PURCHASE_MAX_ACTIVE,
		*constantsBucketData.MAKE_PURCHASE_COOLDOWN,
		*constantsBucketData.CANCEL_PURCHASE_COOLDOWN,

		builderenv.ESI_USER_AGENT,
		builderenv.ESI_MARKETS_CLIENT_ID,
		builderenv.ESI_MARKETS_CLIENT_SECRET,
		builderenv.ESI_CORP_CLIENT_ID,
		builderenv.ESI_CORP_CLIENT_SECRET,
		builderenv.ESI_STRUCTURE_INFO_CLIENT_ID,
		builderenv.ESI_STRUCTURE_INFO_CLIENT_SECRET,
		builderenv.ESI_AUTH_CLIENT_ID,
		builderenv.ESI_AUTH_CLIENT_SECRET,
	)

	_, err = f.WriteString(constantsFileString)
	if err != nil {
		return err
	}

	return nil
}
