package proto

import (
	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
)

type CfgGetMarketNamesParams struct{}

type CfgGetMarketNamesClient struct{}

func NewCfgGetMarketNamesClient() CfgGetMarketNamesClient {
	return CfgGetMarketNamesClient{}
}

func (gmnc CfgGetMarketNamesClient) Fetch(
	x cache.Context,
	params CfgGetMarketNamesParams,
) (
	rep *[]string,
	err error,
) {
	// fetch web shop locations
	webMarkets, _, err := bucket.GetWebMarkets(x)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(webMarkets))
	for name := range webMarkets {
		names = append(names, name)
	}

	return &names, nil
}
