package buildconstants

import "time"

const (
	// // updater data
	CACHE_LOGGING bool = false
	DEV_MODE      bool = false

	// capacities for SDE data
	CAPACITY_SDE_CATEGORIES      int = 0
	CAPACITY_SDE_GROUPS          int = 0
	CAPACITY_SDE_MARKET_GROUPS   int = 0
	CAPACITY_SDE_NAME_TO_TYPE_ID int = 0
	CAPACITY_SDE_REGIONS         int = 0
	CAPACITY_SDE_SYSTEMS         int = 0
	CAPACITY_SDE_SYSTEM_IDS      int = 0
	CAPACITY_SDE_STATIONS        int = 0
	CAPACITY_SDE_TYPE_DATA_MAP   int = 0
	CAPACITY_SDE_TYPE_VOLUMES    int = 0

	// capacities for WEB data
	CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER int = 0
	CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER  int = 0
	CAPACITY_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER     int = 0
	CAPACITY_WEB_BUYBACK_SYSTEMS                  int = 0
	CAPACITY_WEB_SHOP_LOCATIONS                   int = 0
	CAPACITY_WEB_HAUL_ROUTES                      int = 0
	CAPACITY_WEB_MARKETS                          int = 0

	// capacities for CORE data
	CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS int = 0
	CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS  int = 0
	CAPACITY_CORE_HAUL_ROUTE_TYPE_MAPS     int = 0
	CAPACITY_CORE_BUYBACK_SYSTEMS          int = 0
	CAPACITY_CORE_SHOP_LOCATIONS           int = 0
	CAPACITY_CORE_HAUL_ROUTES              int = 0
	CAPACITY_CORE_BANNED_FLAG_SETS         int = 0
	CAPACITY_CORE_HAUL_ROUTE_INFOS         int = 0
	CAPACITY_CORE_PRICINGS                 int = 0
	CAPACITY_CORE_MARKETS                  int = 0

	// // build flags

	// data version
	PROGRAM_VERSION string = "" // probably just git commit hash
	DATA_VERSION    string = "" // bucket data version

	// cache config
	CCACHE_MAX_BYTES int    = 0
	SCACHE_ADDRESS   string = ""

	// needed to bootstrap authentication
	// BOOTSTRAP_ADMIN_ID - only admins can add admins, and there are zero initially
	BOOTSTRAP_ADMIN_ID               int32  = 0
	CORPORATION_ID                   int32  = 0
	CORPORATION_WEB_REFRESH_TOKEN    string = ""
	STRUCTURE_INFO_WEB_REFRESH_TOKEN string = ""

	// RemoteDB Configure
	REMOTEDB            RemoteDB = RDBFirestore
	RDB_MYSQL_HOST      string   = ""
	REMOTEDB_PROJECT_ID string   = ""
	REMOTEDB_CREDS_JSON string   = ""

	// GCP client constructor data
	BUCKET_NAMESPACE  string = ""
	BUCKET_CREDS_JSON string = ""

	// configuration
	PURCHASE_MAX_ACTIVE      int           = 0
	MAKE_PURCHASE_COOLDOWN   time.Duration = 0
	CANCEL_PURCHASE_COOLDOWN time.Duration = 0

	// // Notifications

	// Discord
	DISCORD_BOT_TOKEN string = ""
	DISCORD_CHANNEL   string = ""

	// Buyback Contracts
	BUYBACK_CONTRACT_NOTIFICATIONS          bool   = false
	BUYBACK_CONTRACT_NOTIFICATIONS_BASE_URL string = ""
	// Shop Contracts
	SHOP_CONTRACT_NOTIFICATIONS          bool   = false
	SHOP_CONTRACT_NOTIFICATIONS_BASE_URL string = ""
	// Haul Contracts
	HAUL_CONTRACT_NOTIFICATIONS          bool   = false
	HAUL_CONTRACT_NOTIFICATIONS_BASE_URL string = ""
	// Purchases
	PURCHASE_NOTIFICATIONS          bool   = false
	PURCHASE_NOTIFICATIONS_BASE_URL string = ""

	// // ESI configuration
	ESI_USER_AGENT string = ""
	// esi-markets.structure_markets
	ESI_MARKETS_CLIENT_ID     string = ""
	ESI_MARKETS_CLIENT_SECRET string = ""
	// esi-contracts.read_corporation_contracts
	// esi-assets.read_corporation_assets
	ESI_CORP_CLIENT_ID     string = ""
	ESI_CORP_CLIENT_SECRET string = ""
	// esi-universe.read_structures
	ESI_STRUCTURE_INFO_CLIENT_ID     string = ""
	ESI_STRUCTURE_INFO_CLIENT_SECRET string = ""
	//
	ESI_AUTH_CLIENT_ID     string = ""
	ESI_AUTH_CLIENT_SECRET string = ""
)
