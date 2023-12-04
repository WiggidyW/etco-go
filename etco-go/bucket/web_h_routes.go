package bucket

import (
	"fmt"
	"strconv"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_HAUL_ROUTES_BUF_CAP    int           = 0
	WEB_HAUL_ROUTES_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebHaulRoutes = cache.RegisterType[map[b.SystemId]b.WebHaulRoute]("webhaulroutes", WEB_HAUL_ROUTES_BUF_CAP)
}

func GetWebHaulRoutes(
	x cache.Context,
) (
	rep map[b.WebHaulRouteSystemsKey]b.WebHaulRoute,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebHaulRoutes,
		keys.CacheKeyWebHaulRoutes,
		keys.TypeStrWebHaulRoutes,
		WEB_HAUL_ROUTES_EXPIRES_IN,
		build.CAPACITY_WEB_HAUL_ROUTES,
	)
}

func ProtoGetWebHaulRoutes(
	x cache.Context,
) (
	rep map[string]*proto.CfgHaulRoute,
	expires time.Time,
	err error,
) {
	var webHaulRoutes map[b.WebHaulRouteSystemsKey]b.WebHaulRoute
	webHaulRoutes, expires, err = GetWebHaulRoutes(x)
	if err == nil {
		rep = WebHaulRoutesToProto(webHaulRoutes)
	}
	return rep, expires, err
}

func SetWebHaulRoutes(
	x cache.Context,
	rep map[b.WebHaulRouteSystemsKey]b.WebHaulRoute,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebHaulRoutes,
		keys.CacheKeyWebHaulRoutes,
		keys.TypeStrWebHaulRoutes,
		WEB_HAUL_ROUTES_EXPIRES_IN,
		rep,
		nil,
	)
}

func ProtoMergeSetWebHaulRoutes(
	x cache.Context,
	updates map[string]*proto.CfgHaulRoute,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}
	return protoMergeSetTerritories(
		x,
		updates,
		GetWebBuybackBundleKeys,
		GetWebHaulRoutes,
		ProtoMergeHaulRoutes,
		SetWebHaulRoutes,
	)
}

// // To Proto

func WebHaulRoutesToProto(
	webHaulRoutes map[b.WebHaulRouteSystemsKey]b.WebHaulRoute,
) (
	pbHaulRoutes map[string]*proto.CfgHaulRoute,
) {
	return newPBCfgHaulRoutes(webHaulRoutes)
}

func newPBCfgHaulRoutes(
	webHaulRoutes map[b.WebHaulRouteSystemsKey]b.WebHaulRoute,
) (
	pbHaulRoutes map[string]*proto.CfgHaulRoute,
) {
	pbHaulRoutes = make(
		map[string]*proto.CfgHaulRoute,
		len(webHaulRoutes),
	)
	for systemsKey, webHaulRoute := range webHaulRoutes {
		// unsafe lets us keep the bits as-is. We don't care about endianness.
		_, _, strSystemsKey := WebHaulKeyToStr(systemsKey)
		pbHaulRoutes[strSystemsKey] = newPBCfgHaulRoute(webHaulRoute)
	}
	return pbHaulRoutes
}

func newPBCfgHaulRoute(
	webHaulRoute b.WebHaulRoute,
) (
	pbHaulRoute *proto.CfgHaulRoute,
) {
	return &proto.CfgHaulRoute{
		BundleKey:      webHaulRoute.BundleKey,
		MaxVolume:      webHaulRoute.MaxVolume().Uint32(),
		MinReward:      webHaulRoute.MinReward().Uint64(),
		MaxReward:      webHaulRoute.MaxReward().Uint64(),
		TaxRate:        uint32(webHaulRoute.TaxRate),
		M3Fee:          uint32(webHaulRoute.M3Fee),
		CollateralRate: uint32(webHaulRoute.CollateralRate),
		RewardStrategy: rewardStrategyToProto(webHaulRoute.RewardStrategy),
	}
}

func rewardStrategyToProto(
	rRewardStrategy b.HaulRouteRewardStrategy,
) (
	pbRewardStrategy proto.HaulRewardStrategy,
) {
	switch rRewardStrategy {
	case b.HRRSLesserOf:
		return proto.HaulRewardStrategy_HRS_LESSER_OF
	case b.HRRSGreaterOf:
		return proto.HaulRewardStrategy_HRS_GREATER_OF
	case b.HRRSSum:
		return proto.HaulRewardStrategy_HRS_SUM
	default:
		return proto.HaulRewardStrategy_HRS_INVALID
	}
}

// // Merge

func ProtoMergeHaulRoutes(
	original map[b.WebHaulRouteSystemsKey]b.WebHaulRoute,
	updates map[string]*proto.CfgHaulRoute,
	bundleKeys map[string]struct{},
) (
	err error,
) {
	for strSystemsKey, pbHaulRoute := range updates {
		var startSID, endSID b.SystemId
		var systemsKey b.WebHaulRouteSystemsKey
		startSID, endSID, systemsKey, err = StrToWebHaulKey(strSystemsKey)
		if err != nil {
			return newPBtoWebHaulRouteError(
				systemsKey,
				fmt.Sprintf(
					"invalid key '%s': %s",
					strSystemsKey,
					err,
				),
			)
		} else if pbHaulRoute == nil || pbHaulRoute.BundleKey == "" {
			delete(original, systemsKey)
		} else if _, ok := bundleKeys[pbHaulRoute.BundleKey]; !ok {
			return newPBtoWebHaulRouteError(
				systemsKey,
				fmt.Sprintf(
					"type map key '%s' does not exist",
					pbHaulRoute.BundleKey,
				),
			)
		} else if system := staticdb.GetSystemInfo(startSID); system == nil {
			return newPBtoWebHaulRouteError(
				systemsKey,
				fmt.Sprintf(
					"start system '%d' does not exist",
					startSID,
				),
			)
		} else if system = staticdb.GetSystemInfo(endSID); system == nil {
			return newPBtoWebHaulRouteError(
				systemsKey,
				fmt.Sprintf(
					"end system '%d' does not exist",
					endSID,
				),
			)
		} else {
			original[systemsKey] = pBtoWebHaulRoute(pbHaulRoute)
		}
	}
	return nil
}

func newPBtoWebHaulRouteError(
	systemsKey b.WebHaulRouteSystemsKey,
	errStr string,
) configerror.ErrInvalid {
	startSID, endSID := b.BytesToInt32Pair(systemsKey)
	return configerror.ErrInvalid{
		Err: configerror.ErrHaulRouteInvalid{
			Err: fmt.Errorf(
				"['startSystem: %d', 'endSystem: %d']: %s",
				startSID,
				endSID,
				errStr,
			),
		},
	}
}

func pBtoWebHaulRoute(
	pbHaulRoute *proto.CfgHaulRoute,
) (
	webHaulRoute b.WebHaulRoute,
) {
	maxVolume := b.NewScientific16FromUInt(pbHaulRoute.MaxVolume)
	minReward := b.NewScientific16FromUInt(pbHaulRoute.MinReward)
	maxReward := b.NewScientific16FromUInt(pbHaulRoute.MaxReward)
	return b.WebHaulRoute{
		BundleKey:          pbHaulRoute.BundleKey,
		MaxVolumeS16Base:   maxVolume.Base,
		MaxVolumeS16Zeroes: maxVolume.Zeroes,
		MinRewardS16Base:   minReward.Base,
		MinRewardS16Zeroes: minReward.Zeroes,
		MaxRewardS16Base:   maxReward.Base,
		MaxRewardS16Zeroes: maxReward.Zeroes,
		TaxRate:            b.NewDecPercentage(pbHaulRoute.TaxRate),
		M3Fee:              uint16(pbHaulRoute.M3Fee),
		CollateralRate:     b.NewDecPercentage(pbHaulRoute.CollateralRate),
	}
}

// string - 8bytes conversions

func StrToWebHaulKey(
	strSystemsKey string,
) (
	startSystem b.SystemId,
	endSystem b.SystemId,
	systemsKey b.WebHaulRouteSystemsKey,
	err error,
) {
	if len(strSystemsKey) != 16 {
		return 0, 0, systemsKey, fmt.Errorf(
			"invalid length '%d' for systemsKey '%s'",
			len(strSystemsKey),
			strSystemsKey,
		)
	}
	strStartSystem := strSystemsKey[0:8]
	strEndSystem := strSystemsKey[8:16]
	startSystem, err = hexStrToSystemId(strStartSystem)
	if err != nil {
		return 0, 0, systemsKey, fmt.Errorf(
			"invalid startSystem '%s': %s",
			strStartSystem,
			err,
		)
	}
	endSystem, err = hexStrToSystemId(strEndSystem)
	if err != nil {
		return 0, 0, systemsKey, fmt.Errorf(
			"invalid endSystem '%s': %s",
			strEndSystem,
			err,
		)
	}
	systemsKey = b.Int32PairToBytes(startSystem, endSystem)
	return startSystem, endSystem, systemsKey, nil
}

func WebHaulKeyToStr(
	systemsKey b.WebHaulRouteSystemsKey,
) (
	startSystem b.SystemId,
	endSystem b.SystemId,
	strSystemsKey string,
) {
	startSystem, endSystem = b.BytesToInt32Pair(systemsKey)
	strStartSystem := systemIdToHexStr(startSystem)
	strEndSystem := systemIdToHexStr(endSystem)
	strSystemsKey = strStartSystem + strEndSystem
	return startSystem, endSystem, strSystemsKey
}

func systemIdToHexStr(
	systemId b.SystemId,
) (
	hex string,
) {
	return fmt.Sprintf("%08x", systemId)
}

func hexStrToSystemId(
	hex string,
) (
	systemId b.SystemId,
	err error,
) {
	var systemId64 int64
	systemId64, err = strconv.ParseInt(hex, 16, 32)
	systemId = b.SystemId(systemId64)
	return systemId, err
}
