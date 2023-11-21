package esi

import (
	"fmt"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

type EsiApp uint8

func AppFromProto(protoApp proto.EsiApp) (
	app EsiApp,
	err error,
) {
	switch protoApp {
	case proto.EsiApp_EA_AUTH:
		return EsiAppAuth, nil
	case proto.EsiApp_EA_CORPORATION:
		return EsiAppCorp, nil
	case proto.EsiApp_EA_STRUCTURE_INFO:
		return EsiAppStructureInfo, nil
	case proto.EsiApp_EA_MARKETS:
		return EsiAppMarkets, nil
	default:
		return EsiAppProtoInvalid, protoerr.MsgNew(
			protoerr.INVALID_REQUEST,
			fmt.Sprintf("Invalid EsiApp: %d", protoApp),
		)
	}
}

const (
	EsiAppAuth EsiApp = iota
	EsiAppCorp
	EsiAppStructureInfo
	EsiAppMarkets
	EsiAppProtoInvalid
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
