package builder

import (
	"fmt"
	"os"
	"strconv"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go-builder/builderenv"
)

func transceiveWriteConstants(
	filePath string,
	constantsData b.ConstantsData,
	sdeUpdaterData b.SDEUpdaterData,
	coreUpdaterData b.CoreUpdaterData,
	chnSendDone chanresult.ChanSendResult[struct{}],
) error {
	err := writeConstants(
		filePath,
		constantsData,
		sdeUpdaterData,
		coreUpdaterData,
	)
	if err != nil {
		return chnSendDone.SendErr(err)
	} else {
		return chnSendDone.SendOk(struct{}{})
	}
}

func writeConstants(
	filePath string,
	constantsData b.ConstantsData,
	sdeUpdaterData b.SDEUpdaterData,
	coreUpdaterData b.CoreUpdaterData,
) error {
	constantsData = useEnvAndDefaultsIfNil(constantsData)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	constantsFileString := fmt.Sprintf(
		`package buildconstants

		import "time"

		const (
			CACHE_LOGGING bool = %s
			DEV_MODE      bool = %s

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

			PROGRAM_VERSION string = "%s"
			DATA_VERSION    string = "%s"

			CCACHE_MAX_BYTES int    = %d
			SCACHE_ADDRESS   string = "%s"

			BOOTSTRAP_ADMIN_ID               int32  = %d
			CORPORATION_ID                   int32  = %d
			CORPORATION_WEB_REFRESH_TOKEN    string = "%s"
			STRUCTURE_INFO_WEB_REFRESH_TOKEN string = "%s"

			REMOTEDB 		  RemoteDB = %s
			RDB_MYSQL_HOST      string = "%s"
			REMOTEDB_PROJECT_ID string = "%s"
			REMOTEDB_CREDS_JSON string = %s

			BUCKET_NAMESPACE    string = "%s"
			BUCKET_CREDS_JSON   string = %s

			PURCHASE_MAX_ACTIVE      int = %d
			MAKE_PURCHASE_COOLDOWN   time.Duration = %d
			CANCEL_PURCHASE_COOLDOWN time.Duration = %d

			DISCORD_BOT_TOKEN string = "%s"
			DISCORD_CHANNEL   string = "%s"

			BUYBACK_CONTRACT_NOTIFICATIONS          bool   = %s
			BUYBACK_CONTRACT_NOTIFICATIONS_BASE_URL string = "%s"
			SHOP_CONTRACT_NOTIFICATIONS 			bool   = %s
			SHOP_CONTRACT_NOTIFICATIONS_BASE_URL 	string = "%s"
			PURCHASE_NOTIFICATIONS 					bool   = %s
			PURCHASE_NOTIFICATIONS_BASE_URL 		string = "%s"

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
		strconv.FormatBool(builderenv.CACHE_LOGGING),
		strconv.FormatBool(builderenv.DEV_MODE),

		sdeUpdaterData.CAPACITY_SDE_CATEGORIES,
		sdeUpdaterData.CAPACITY_SDE_GROUPS,
		sdeUpdaterData.CAPACITY_SDE_MARKET_GROUPS,
		sdeUpdaterData.CAPACITY_SDE_NAME_TO_TYPE_ID,
		sdeUpdaterData.CAPACITY_SDE_REGIONS,
		sdeUpdaterData.CAPACITY_SDE_SYSTEMS,
		sdeUpdaterData.CAPACITY_SDE_STATIONS,
		sdeUpdaterData.CAPACITY_SDE_TYPE_DATA_MAP,
		sdeUpdaterData.CAPACITY_SDE_TYPE_VOLUMES,

		coreUpdaterData.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
		coreUpdaterData.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
		coreUpdaterData.CAPACITY_WEB_BUYBACK_SYSTEMS,
		coreUpdaterData.CAPACITY_WEB_SHOP_LOCATIONS,
		coreUpdaterData.CAPACITY_WEB_MARKETS,

		coreUpdaterData.CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS,
		coreUpdaterData.CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS,
		coreUpdaterData.CAPACITY_CORE_BUYBACK_SYSTEMS,
		coreUpdaterData.CAPACITY_CORE_SHOP_LOCATIONS,
		coreUpdaterData.CAPACITY_CORE_BANNED_FLAG_SETS,
		coreUpdaterData.CAPACITY_CORE_PRICINGS,
		coreUpdaterData.CAPACITY_CORE_MARKETS,

		builderenv.PROGRAM_VERSION,
		coreUpdaterData.VERSION_BUYBACK,

		builderenv.CCACHE_MAX_BYTES,
		builderenv.SCACHE_ADDRESS,

		builderenv.BOOTSTRAP_ADMIN_ID,
		builderenv.CORPORATION_ID,
		*constantsData.CORPORATION_WEB_REFRESH_TOKEN,
		*constantsData.STRUCTURE_INFO_WEB_REFRESH_TOKEN,

		builderenv.REMOTEDB,
		builderenv.RDB_MYSQL_HOST,
		builderenv.REMOTEDB_PROJECT_ID,
		strconv.Quote(builderenv.REMOTEDB_CREDS_JSON),

		builderenv.BUCKET_NAMESPACE,
		strconv.Quote(builderenv.BUCKET_CREDS_JSON),

		*constantsData.PURCHASE_MAX_ACTIVE,
		*constantsData.MAKE_PURCHASE_COOLDOWN,
		*constantsData.CANCEL_PURCHASE_COOLDOWN,

		builderenv.DISCORD_BOT_TOKEN,
		*constantsData.DISCORD_CHANNEL,

		strconv.FormatBool(*constantsData.BUYBACK_CONTRACT_NOTIFICATIONS),
		builderenv.BUYBACK_CONTRACT_NOTIFICATIONS_BASE_URL,
		strconv.FormatBool(*constantsData.SHOP_CONTRACT_NOTIFICATIONS),
		builderenv.SHOP_CONTRACT_NOTIFICATIONS_BASE_URL,
		strconv.FormatBool(*constantsData.PURCHASE_NOTIFICATIONS),
		builderenv.PURCHASE_NOTIFICATIONS_BASE_URL,

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

func useEnvAndDefaultsIfNil(constantsData b.ConstantsData) b.ConstantsData {
	// If any values are missing from bucket data, set them to ENV values.
	if constantsData.PURCHASE_MAX_ACTIVE == nil {
		constantsData.PURCHASE_MAX_ACTIVE =
			&builderenv.PURCHASE_MAX_ACTIVE
	}
	if constantsData.MAKE_PURCHASE_COOLDOWN == nil {
		constantsData.MAKE_PURCHASE_COOLDOWN =
			&builderenv.MAKE_PURCHASE_COOLDOWN
	}
	if constantsData.CANCEL_PURCHASE_COOLDOWN == nil {
		constantsData.CANCEL_PURCHASE_COOLDOWN =
			&builderenv.CANCEL_PURCHASE_COOLDOWN
	}
	if constantsData.CORPORATION_WEB_REFRESH_TOKEN == nil ||
		*constantsData.CORPORATION_WEB_REFRESH_TOKEN == "" {
		constantsData.CORPORATION_WEB_REFRESH_TOKEN =
			&builderenv.CORPORATION_WEB_REFRESH_TOKEN
	}
	if constantsData.STRUCTURE_INFO_WEB_REFRESH_TOKEN == nil ||
		*constantsData.STRUCTURE_INFO_WEB_REFRESH_TOKEN == "" {
		constantsData.STRUCTURE_INFO_WEB_REFRESH_TOKEN =
			&builderenv.STRUCTURE_INFO_WEB_REFRESH_TOKEN
	}
	if constantsData.DISCORD_CHANNEL == nil ||
		*constantsData.DISCORD_CHANNEL == "" {
		constantsData.DISCORD_CHANNEL = &builderenv.DISCORD_CHANNEL
	}
	if constantsData.BUYBACK_CONTRACT_NOTIFICATIONS == nil {
		constantsData.BUYBACK_CONTRACT_NOTIFICATIONS =
			&builderenv.BUYBACK_CONTRACT_NOTIFICATIONS
	}
	if constantsData.SHOP_CONTRACT_NOTIFICATIONS == nil {
		constantsData.SHOP_CONTRACT_NOTIFICATIONS =
			&builderenv.SHOP_CONTRACT_NOTIFICATIONS
	}
	if constantsData.PURCHASE_NOTIFICATIONS == nil {
		constantsData.PURCHASE_NOTIFICATIONS =
			&builderenv.PURCHASE_NOTIFICATIONS
	}
	return constantsData
}
