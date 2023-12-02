package builderenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_REMOTEDB = "firestore"

var (
	// enum strings
	STR_REMOTEDB = os.Getenv("REMOTEDB")

	// boolean strings ('true' if 'true', 'false' if anything else)
	STR_CACHE_LOGGING                  = os.Getenv("CACHE_LOGGING")
	STR_DEV_MODE                       = os.Getenv("DEV_MODE")
	STR_BUYBACK_CONTRACT_NOTIFICATIONS = os.Getenv("BUYBACK_CONTRACT_NOTIFICATIONS")
	STR_SHOP_CONTRACT_NOTIFICATIONS    = os.Getenv("SHOP_CONTRACT_NOTIFICATIONS")
	STR_HAUL_CONTRACT_NOTIFICATIONS    = os.Getenv("HAUL_CONTRACT_NOTIFICATIONS")
	STR_PURCHASE_NOTIFICATIONS         = os.Getenv("PURCHASE_NOTIFICATIONS")

	// integer strings
	STR_CORPORATION_ID      = os.Getenv("CORPORATION_ID")
	STR_BOOTSTRAP_ADMIN_ID  = os.Getenv("BOOTSTRAP_ADMIN_ID")
	STR_PURCHASE_MAX_ACTIVE = os.Getenv("PURCHASE_MAX_ACTIVE")
	STR_CCACHE_MAX_BYTES    = os.Getenv("CCACHE_MAX_BYTES")

	// duration strings
	STR_MAKE_PURCHASE_COOLDOWN   = os.Getenv("MAKE_PURCHASE_COOLDOWN")
	STR_CANCEL_PURCHASE_COOLDOWN = os.Getenv("CANCEL_PURCHASE_COOLDOWN")

	// actual strings
	RDB_MYSQL_HOST                          = os.Getenv("RDB_MYSQL_HOST")
	PROGRAM_VERSION                         = os.Getenv("PROGRAM_VERSION")
	CORPORATION_WEB_REFRESH_TOKEN           = os.Getenv("CORPORATION_WEB_REFRESH_TOKEN")
	STRUCTURE_INFO_WEB_REFRESH_TOKEN        = os.Getenv("STRUCTURE_INFO_WEB_REFRESH_TOKEN")
	REMOTEDB_PROJECT_ID                     = os.Getenv("REMOTEDB_PROJECT_ID")
	REMOTEDB_CREDS_JSON                     = os.Getenv("REMOTEDB_CREDS_JSON")
	BUCKET_NAMESPACE                        = os.Getenv("BUCKET_NAMESPACE")
	BUCKET_CREDS_JSON                       = os.Getenv("BUCKET_CREDS_JSON")
	ESI_USER_AGENT                          = os.Getenv("ESI_USER_AGENT")
	ESI_MARKETS_CLIENT_ID                   = os.Getenv("ESI_MARKETS_CLIENT_ID")
	ESI_MARKETS_CLIENT_SECRET               = os.Getenv("ESI_MARKETS_CLIENT_SECRET")
	ESI_CORP_CLIENT_ID                      = os.Getenv("ESI_CORP_CLIENT_ID")
	ESI_CORP_CLIENT_SECRET                  = os.Getenv("ESI_CORP_CLIENT_SECRET")
	ESI_STRUCTURE_INFO_CLIENT_ID            = os.Getenv("ESI_STRUCTURE_INFO_CLIENT_ID")
	ESI_STRUCTURE_INFO_CLIENT_SECRET        = os.Getenv("ESI_STRUCTURE_INFO_CLIENT_SECRET")
	ESI_AUTH_CLIENT_ID                      = os.Getenv("ESI_AUTH_CLIENT_ID")
	ESI_AUTH_CLIENT_SECRET                  = os.Getenv("ESI_AUTH_CLIENT_SECRET")
	SCACHE_ADDRESS                          = os.Getenv("SCACHE_ADDRESS")
	DISCORD_BOT_TOKEN                       = os.Getenv("DISCORD_BOT_TOKEN")
	DISCORD_CHANNEL                         = os.Getenv("DISCORD_CHANNEL")
	BUYBACK_CONTRACT_NOTIFICATIONS_BASE_URL = os.Getenv("BUYBACK_CONTRACT_NOTIFICATIONS_BASE_URL")
	SHOP_CONTRACT_NOTIFICATIONS_BASE_URL    = os.Getenv("SHOP_CONTRACT_NOTIFICATIONS_BASE_URL")
	HAUL_CONTRACT_NOTIFICATIONS_BASE_URL    = os.Getenv("HAUL_CONTRACT_NOTIFICATIONS_BASE_URL")
	PURCHASE_NOTIFICATIONS_BASE_URL         = os.Getenv("PURCHASE_NOTIFICATIONS_BASE_URL")

	// file paths
	GOB_FILE_DIR        = os.Getenv("GOB_FILE_DIR")
	CONSTANTS_FILE_PATH = os.Getenv("CONSTANTS_FILE_PATH")

	// booleans
	DEV_MODE                       bool = false
	BUYBACK_CONTRACT_NOTIFICATIONS bool = false
	SHOP_CONTRACT_NOTIFICATIONS    bool = false
	HAUL_CONTRACT_NOTIFICATIONS    bool = false
	PURCHASE_NOTIFICATIONS         bool = false
	CACHE_LOGGING                  bool = false

	// enums (string form)
	REMOTEDB string = ""
)

var (
	CORPORATION_ID           int32         = 0
	BOOTSTRAP_ADMIN_ID       int32         = 0
	PURCHASE_MAX_ACTIVE      int           = 0
	MAKE_PURCHASE_COOLDOWN   time.Duration = 0
	CANCEL_PURCHASE_COOLDOWN time.Duration = 0
	CCACHE_MAX_BYTES         int           = 0
)

func ConvertAndValidate() (err error) {
	STR_REMOTEDB = strings.ToLower(STR_REMOTEDB)
	STR_DEV_MODE = strings.ToLower(STR_DEV_MODE)
	STR_BUYBACK_CONTRACT_NOTIFICATIONS = strings.ToLower(
		STR_BUYBACK_CONTRACT_NOTIFICATIONS,
	)
	STR_SHOP_CONTRACT_NOTIFICATIONS = strings.ToLower(
		STR_SHOP_CONTRACT_NOTIFICATIONS,
	)
	STR_HAUL_CONTRACT_NOTIFICATIONS = strings.ToLower(
		STR_HAUL_CONTRACT_NOTIFICATIONS,
	)
	STR_PURCHASE_NOTIFICATIONS = strings.ToLower(STR_PURCHASE_NOTIFICATIONS)
	STR_CACHE_LOGGING = strings.ToLower(STR_CACHE_LOGGING)

	if STR_DEV_MODE != "true" {
		STR_DEV_MODE = "false"
		DEV_MODE = false
	} else {
		DEV_MODE = true
	}
	if STR_BUYBACK_CONTRACT_NOTIFICATIONS != "true" {
		STR_BUYBACK_CONTRACT_NOTIFICATIONS = "false"
		BUYBACK_CONTRACT_NOTIFICATIONS = false
	} else {
		BUYBACK_CONTRACT_NOTIFICATIONS = true
	}
	if STR_SHOP_CONTRACT_NOTIFICATIONS != "true" {
		STR_SHOP_CONTRACT_NOTIFICATIONS = "false"
		SHOP_CONTRACT_NOTIFICATIONS = false
	} else {
		SHOP_CONTRACT_NOTIFICATIONS = true
	}
	if STR_HAUL_CONTRACT_NOTIFICATIONS != "true" {
		STR_HAUL_CONTRACT_NOTIFICATIONS = "false"
		HAUL_CONTRACT_NOTIFICATIONS = false
	} else {
		HAUL_CONTRACT_NOTIFICATIONS = true
	}
	if STR_PURCHASE_NOTIFICATIONS != "true" {
		STR_PURCHASE_NOTIFICATIONS = "false"
		PURCHASE_NOTIFICATIONS = false
	} else {
		PURCHASE_NOTIFICATIONS = true
	}
	if STR_CACHE_LOGGING != "true" {
		STR_CACHE_LOGGING = "false"
		CACHE_LOGGING = false
	} else {
		CACHE_LOGGING = true
	}

	// validate RemoteDB Config
	if STR_REMOTEDB != "mysql" && STR_REMOTEDB != "firestore" {
		STR_REMOTEDB = DEFAULT_REMOTEDB
	}
	if STR_REMOTEDB == "mysql" {
		REMOTEDB = "RDBMySQL"
		if RDB_MYSQL_HOST == "" {
			return fmt.Errorf("REMOTEDB is 'mysql' and RDB_MYSQL_HOST is empty")
		}
	} else if STR_REMOTEDB == "firestore" {
		REMOTEDB = "RDBFirestore"
		if REMOTEDB_PROJECT_ID == "" {
			return fmt.Errorf(
				"REMOTEDB is 'firestore' and REMOTEDB_PROJECT_ID is empty",
			)
		} else if REMOTEDB_CREDS_JSON == "" {
			return fmt.Errorf(
				"REMOTEDB is 'firestore' and REMOTEDB_CREDS_JSON is empty",
			)
		}
	}

	// ensure that no env vars are empty or missing
	if STR_CORPORATION_ID == "" {
		return fmt.Errorf("CORPORATION_ID is empty")
	} else if STR_BOOTSTRAP_ADMIN_ID == "" {
		return fmt.Errorf("BOOTSTRAP_ADMIN_ID is empty")
	} else if CORPORATION_WEB_REFRESH_TOKEN == "" {
		return fmt.Errorf("CORPORATION_WEB_REFRESH_TOKEN is empty")
	} else if STRUCTURE_INFO_WEB_REFRESH_TOKEN == "" {
		return fmt.Errorf("STRUCTURE_INFO_WEB_REFRESH_TOKEN is empty")
	} else if BUCKET_NAMESPACE == "" {
		return fmt.Errorf("BUCKET_NAMESPACE is empty")
	} else if BUCKET_CREDS_JSON == "" {
		return fmt.Errorf("BUCKET_CREDS_JSON is empty")
	} else if ESI_USER_AGENT == "" {
		return fmt.Errorf("ESI_USER_AGENT is empty")
	} else if ESI_MARKETS_CLIENT_ID == "" {
		return fmt.Errorf("ESI_MARKETS_CLIENT_ID is empty")
	} else if ESI_MARKETS_CLIENT_SECRET == "" {
		return fmt.Errorf("ESI_MARKETS_CLIENT_SECRET is empty")
	} else if ESI_CORP_CLIENT_ID == "" {
		return fmt.Errorf("ESI_CORP_CLIENT_ID is empty")
	} else if ESI_CORP_CLIENT_SECRET == "" {
		return fmt.Errorf("ESI_CORP_CLIENT_SECRET is empty")
	} else if ESI_STRUCTURE_INFO_CLIENT_ID == "" {
		return fmt.Errorf("ESI_STRUCTURE_INFO_CLIENT_ID is empty")
	} else if ESI_STRUCTURE_INFO_CLIENT_SECRET == "" {
		return fmt.Errorf("ESI_STRUCTURE_INFO_CLIENT_SECRET is empty")
	} else if ESI_AUTH_CLIENT_ID == "" {
		return fmt.Errorf("ESI_AUTH_CLIENT_ID is empty")
	} else if ESI_AUTH_CLIENT_SECRET == "" {
		return fmt.Errorf("ESI_AUTH_CLIENT_SECRET is empty")
	} else if STR_PURCHASE_MAX_ACTIVE == "" {
		return fmt.Errorf("PURCHASE_MAX_ACTIVE is empty")
	} else if STR_MAKE_PURCHASE_COOLDOWN == "" {
		return fmt.Errorf("MAKE_PURCHASE_COOLDOWN is empty")
	} else if STR_CANCEL_PURCHASE_COOLDOWN == "" {
		return fmt.Errorf("CANCEL_PURCHASE_COOLDOWN is empty")
	} else if GOB_FILE_DIR == "" {
		return fmt.Errorf("GOB_FILE_DIR is empty")
	} else if CONSTANTS_FILE_PATH == "" {
		return fmt.Errorf("CONSTANTS_FILE_PATH is empty")
	} else if STR_CCACHE_MAX_BYTES == "" {
		return fmt.Errorf("CCACHE_MAX_BYTES is empty")
	} else if SCACHE_ADDRESS == "" {
		return fmt.Errorf("SCACHE_ADDRESS is empty")
	} else if PROGRAM_VERSION == "" {
		return fmt.Errorf("PROGRAM_VERSION is empty")
	}

	// ensure that the string ints are valid and convert them
	if I64_CORPORATION_ID, err := strconv.ParseInt(
		STR_CORPORATION_ID,
		10,
		32,
	); err != nil {
		return err
	} else {
		CORPORATION_ID = int32(I64_CORPORATION_ID)
	}
	if I64_BOOTSTRAP_ADMIN_ID, err := strconv.ParseInt(
		STR_BOOTSTRAP_ADMIN_ID,
		10,
		32,
	); err != nil {
		return err
	} else {
		BOOTSTRAP_ADMIN_ID = int32(I64_BOOTSTRAP_ADMIN_ID)
	}
	if I64_PURCHASE_MAX_ACTIVE, err := strconv.ParseInt(
		STR_PURCHASE_MAX_ACTIVE,
		10,
		0,
	); err != nil {
		return err
	} else {
		PURCHASE_MAX_ACTIVE = int(I64_PURCHASE_MAX_ACTIVE)
	}
	if I64_CCACHE_MAX_BYTES, err := strconv.ParseInt(
		STR_CCACHE_MAX_BYTES,
		10,
		0,
	); err != nil {
		return err
	} else {
		CCACHE_MAX_BYTES = int(I64_CCACHE_MAX_BYTES)
	}

	// ensure that the string durations are valid and convert them
	if MAKE_PURCHASE_COOLDOWN, err = time.ParseDuration(
		STR_MAKE_PURCHASE_COOLDOWN,
	); err != nil {
		return err
	}

	if CANCEL_PURCHASE_COOLDOWN, err = time.ParseDuration(
		STR_CANCEL_PURCHASE_COOLDOWN,
	); err != nil {
		return err
	}

	// validate GOB_FILE_DIR
	if err := validateCreateDirectory(GOB_FILE_DIR); err != nil {
		return err
	}

	// validate CONSTANTS_FILE_PATH
	if err := validateFileAndValidateCreateFileDirectory(
		CONSTANTS_FILE_PATH,
	); err != nil {
		return err
	}

	return nil
}
