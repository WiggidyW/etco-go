package webtocore

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

func convertWebMarkets(webMarkets map[b.MarketName]b.WebMarket) (
	coreMarkets []b.Market,
	coreMarketsIndexMap map[b.MarketName]int,
) {
	coreMarkets = make([]b.Market, 0, len(webMarkets))
	coreMarketsIndexMap = make(map[b.MarketName]int, len(webMarkets))

	for marketName, webMarket := range webMarkets {
		coreMarketsIndexMap[marketName] = len(coreMarkets)
		coreMarkets = append(coreMarkets, b.Market{
			Name:         marketName,
			RefreshToken: webMarket.RefreshToken,
			LocationId:   webMarket.LocationId,
			IsStructure:  webMarket.IsStructure,
		})
	}

	return coreMarkets, coreMarketsIndexMap
}
