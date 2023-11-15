package esi

import (
	"fmt"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/logger"
)

type EsiApp uint8

const (
	EsiAppAuth EsiApp = iota
	EsiAppCorp
	EsiAppStructureInfo
	EsiAppMarkets
)

func (app EsiApp) String() string {
	switch app {
	case EsiAppAuth:
		return "Auth"
	case EsiAppCorp:
		return "Corp"
	case EsiAppStructureInfo:
		return "StructureInfo"
	case EsiAppMarkets:
		return "Markets"
	default:
		logger.Fatal(fmt.Sprintf("Unknown EsiApp: %d", app))
		return ""
	}
}

func (app EsiApp) TypeStrToken() string {
	switch app {
	case EsiAppAuth:
		return keys.TypeStrAuthToken
	case EsiAppCorp:
		return keys.TypeStrCorpToken
	case EsiAppStructureInfo:
		return keys.TypeStrStructureInfoToken
	case EsiAppMarkets:
		return keys.TypeStrMarketsToken
	default:
		logger.Fatal(fmt.Sprintf("Unknown EsiApp: %d", app))
		return ""
	}
}

func (app EsiApp) CacheKeyToken(token string) string {
	switch app {
	case EsiAppAuth:
		return keys.CacheKeyAuthToken(token)
	case EsiAppCorp:
		return keys.CacheKeyCorpToken(token)
	case EsiAppStructureInfo:
		return keys.CacheKeyStructureInfoToken(token)
	case EsiAppMarkets:
		return keys.CacheKeyMarketsToken(token)
	default:
		logger.Fatal(fmt.Sprintf("Unknown EsiApp: %d", app))
		return ""
	}
}

func (app EsiApp) ClientIdAndSecret() (
	clientId string,
	clientSecret string,
) {
	switch app {
	case EsiAppAuth:
		return build.ESI_AUTH_CLIENT_ID, build.ESI_AUTH_CLIENT_SECRET
	case EsiAppCorp:
		return build.ESI_CORP_CLIENT_ID, build.ESI_CORP_CLIENT_SECRET
	case EsiAppStructureInfo:
		return build.ESI_STRUCTURE_INFO_CLIENT_ID, build.ESI_STRUCTURE_INFO_CLIENT_SECRET
	case EsiAppMarkets:
		return build.ESI_MARKETS_CLIENT_ID, build.ESI_MARKETS_CLIENT_SECRET
	default:
		logger.Fatal(fmt.Sprintf("Unknown EsiApp: %d", app))
		return "", ""
	}
}
