package webtocore

import (
	"errors"
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

func convertWebHaulRoutes(
	webHaulRoutes map[b.WebHaulRouteSystemsKey]b.WebHaulRoute,
	coreHRTypeMapsIndexMap map[b.BundleKey]int,
	sdeSystems map[b.SystemId]b.System,
) (
	coreHaulRoutes map[b.HaulRouteSystemsKey]b.HaulRouteInfoIndex,
	coreHaulRouteInfos []b.HaulRouteInfo,
	err error,
) {
	coreHaulRoutes = make(
		map[b.HaulRouteSystemsKey]b.HaulRouteInfoIndex,
		len(webHaulRoutes),
	)
	coreHaulRouteInfos = make([]b.HaulRouteInfo, 0)
	coreHaulRouteInfosIndexMap := make(map[[14]byte]b.HaulRouteInfoIndex)

	for systemsKey, webHaulRoute := range webHaulRoutes {
		typeMapIndex, ok := coreHRTypeMapsIndexMap[webHaulRoute.BundleKey]
		if !ok {
			return nil, nil, fmt.Errorf(
				"HaulRoute %s has invalid BundleKey %s",
				systemsKey,
				webHaulRoute.BundleKey,
			)
		} else if typeMapIndex > 255 {
			return nil, nil, errors.New(
				"Only a maximum of 255 HaulRouteTypeMaps are supported",
			)
		}
		haulRouteInfoIndex := addHaulRouteInfo(
			&coreHaulRouteInfos,
			coreHaulRouteInfosIndexMap,
			webHaulRoute,
			uint8(typeMapIndex),
		)
		err = addHaulRoute(
			coreHaulRoutes,
			systemsKey,
			haulRouteInfoIndex,
			sdeSystems,
		)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(coreHaulRouteInfos) > 65535 {
		return nil, nil, errors.New(
			"Only a maximum of 65535 HaulRouteInfos are supported",
		)
	}

	return coreHaulRoutes, coreHaulRouteInfos, nil
}

func addHaulRoute(
	haulRoutes map[b.HaulRouteSystemsKey]b.HaulRouteInfoIndex,
	webSystemsKey b.WebHaulRouteSystemsKey,
	haulRouteInfoIndex b.HaulRouteInfoIndex,
	sdeSystems map[b.SystemId]b.System,
) (err error) {
	startSystemId, endSystemId := b.BytesToInt32Pair(webSystemsKey)
	startSystem, startOk := sdeSystems[startSystemId]
	endSystem, endOk := sdeSystems[endSystemId]
	if !startOk || !endOk {
		return fmt.Errorf(
			"HaulRoute has invalid systems: 'start: %d', 'end: %d'",
			startSystemId,
			endSystemId,
		)
	}
	systemsKey := b.Uint16PairToBytes(startSystem.Index, endSystem.Index)
	haulRoutes[systemsKey] = haulRouteInfoIndex
	return nil
}

func addHaulRouteInfo(
	haulRouteInfos *[]b.HaulRouteInfo,
	haulRouteInfoIndexMap map[[14]byte]b.HaulRouteInfoIndex,
	webHaulRoute b.WebHaulRoute,
	typeMapindex uint8,
) (haulRouteInfoIndex b.HaulRouteInfoIndex) {
	taxRateBytes := b.Uint16ToBytes(webHaulRoute.TaxRate)
	m3FeeBytes := b.Uint16ToBytes(webHaulRoute.M3Fee)
	collateralRateBytes := b.Uint16ToBytes(webHaulRoute.CollateralRate)
	infoHashKey := [14]byte{
		webHaulRoute.MaxVolumeS16Base,
		webHaulRoute.MaxVolumeS16Zeroes,
		webHaulRoute.MinRewardS16Base,
		webHaulRoute.MinRewardS16Zeroes,
		webHaulRoute.MaxRewardS16Base,
		webHaulRoute.MaxRewardS16Zeroes,
		taxRateBytes[0],
		taxRateBytes[1],
		m3FeeBytes[0],
		m3FeeBytes[1],
		collateralRateBytes[0],
		collateralRateBytes[1],
		uint8(webHaulRoute.RewardStrategy),
		typeMapindex,
	}

	var ok bool
	haulRouteInfoIndex, ok = haulRouteInfoIndexMap[infoHashKey]
	if !ok {
		haulRouteInfoIndex = uint16(len(*haulRouteInfos))
		*haulRouteInfos = append(*haulRouteInfos, b.HaulRouteInfo{
			MaxVolumeS16Base:   webHaulRoute.MaxVolumeS16Base,
			MaxVolumeS16Zeroes: webHaulRoute.MaxVolumeS16Zeroes,
			MinRewardS16Base:   webHaulRoute.MinRewardS16Base,
			MinRewardS16Zeroes: webHaulRoute.MinRewardS16Zeroes,
			MaxRewardS16Base:   webHaulRoute.MaxRewardS16Base,
			MaxRewardS16Zeroes: webHaulRoute.MaxRewardS16Zeroes,
			TaxRate:            webHaulRoute.TaxRate,
			M3Fee:              webHaulRoute.M3Fee,
			CollateralRate:     webHaulRoute.CollateralRate,
			RewardStrategy:     webHaulRoute.RewardStrategy,
			TypeMapIndex:       typeMapindex,
		})
	}

	return haulRouteInfoIndex
}
