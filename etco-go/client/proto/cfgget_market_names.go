package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/bucket"
)

type CfgGetMarketNamesParams struct{}

type CfgGetMarketNamesClient struct {
	webMarketsReaderClient bucket.SC_WebMarketsReaderClient
}

func NewCfgGetMarketNamesClient(
	webMarketsReaderClient bucket.SC_WebMarketsReaderClient,
) CfgGetMarketNamesClient {
	return CfgGetMarketNamesClient{webMarketsReaderClient}
}

func (gmnc CfgGetMarketNamesClient) Fetch(
	ctx context.Context,
	params CfgGetMarketNamesParams,
) (
	rep *[]string,
	err error,
) {
	// fetch web shop locations
	webMarkets, err := gmnc.webMarketsReaderClient.Fetch(
		ctx,
		bucket.WebMarketsReaderParams{},
	)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(webMarkets.Data()))
	for name := range webMarkets.Data() {
		names = append(names, name)
	}

	return &names, nil
}
