package buildconstants

import "time"

const (
	// // updater data

	// capacities for SDE data
	CAPACITY_SDE_CATEGORIES      int = 0
	CAPACITY_SDE_GROUPS          int = 0
	CAPACITY_SDE_MARKET_GROUPS   int = 0
	CAPACITY_SDE_NAME_TO_TYPE_ID int = 0
	CAPACITY_SDE_REGIONS         int = 0
	CAPACITY_SDE_SYSTEMS         int = 0
	CAPACITY_SDE_STATIONS        int = 0
	CAPACITY_SDE_TYPE_DATA_MAP   int = 0
	CAPACITY_SDE_TYPE_VOLUMES    int = 0

	// capacities for WEB data
	CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER int = 0
	CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER  int = 0
	CAPACITY_WEB_MARKETS                          int = 0
	CAPACITY_WEB_SHOP_LOCATIONS                   int = 0
	CAPACITY_WEB_BUYBACK_SYSTEMS                  int = 0

	// capacities for CORE data
	CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS  int = 0
	CAPACITY_CORE_SHOP_LOCATIONS           int = 0
	CAPACITY_CORE_BANNED_FLAG_SETS         int = 0
	CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS int = 0
	CAPACITY_CORE_BUYBACK_SYSTEMS          int = 0
	CAPACITY_CORE_PRICINGS                 int = 0
	CAPACITY_CORE_MARKETS                  int = 0

	// buyback and shop versions (== updated time of respective buckets)
	VERSION_BUYBACK string = ""
	VERSION_SHOP    string = ""

	// // build flags

	// needed to bootstrap authentication
	// only admins can add admins, and there are zero initially
	BOOTSTRAP_ADMIN_ID               int32  = 0
	CORPORATION_ID                   int32  = 0
	CORPORATION_WEB_REFRESH_TOKEN    string = ""
	STRUCTURE_INFO_WEB_REFRESH_TOKEN string = ""

	// GCP client constructor data
	REMOTEDB_PROJECT_ID string = ""
	REMOTEDB_CREDS_JSON string = ""
	BUCKET_CREDS_JSON   string = ""

	// configuration
	PURCHASE_MAX_ACTIVE      int           = 0
	MAKE_PURCHASE_COOLDOWN   time.Duration = 0
	CANCEL_PURCHASE_COOLDOWN time.Duration = 0

	// ESI configuration
	ESI_USER_AGENT                   string = ""
	ESI_MARKETS_CLIENT_ID            string = ""
	ESI_MARKETS_CLIENT_SECRET        string = ""
	ESI_CORP_CLIENT_ID               string = ""
	ESI_CORP_CLIENT_SECRET           string = ""
	ESI_STRUCTURE_INFO_CLIENT_ID     string = ""
	ESI_STRUCTURE_INFO_CLIENT_SECRET string = ""
	ESI_AUTH_CLIENT_ID               string = ""
	ESI_AUTH_CLIENT_SECRET           string = ""
)
