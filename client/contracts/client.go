package contracts

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
	a "github.com/WiggidyW/eve-trading-co-go/client/authingfwding/authing"
	"github.com/WiggidyW/eve-trading-co-go/client/caching"
	wc "github.com/WiggidyW/eve-trading-co-go/client/caching/weak"
	ct "github.com/WiggidyW/eve-trading-co-go/client/esi/model/contractscorporation"
	"github.com/WiggidyW/eve-trading-co-go/staticdb"
)

type A_WC_ContractsClient = a.AuthingClient[
	authingfwding.WithAuthableParams[ContractsParams],
	ContractsParams,
	caching.CachingResponse[Contracts],
	WC_ContractsClient,
]

type WC_ContractsClient = wc.WeakCachingClient[
	ContractsParams,
	Contracts,
	cache.ExpirableData[Contracts],
	ContractsClient,
]

type ContractsClient struct {
	Inner ct.ContractsCorporationClient
}

func (cc ContractsClient) Fetch(
	ctx context.Context,
	params ContractsParams,
) (*cache.ExpirableData[Contracts], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := cc.Inner.Fetch(ctx, ct.ContractsCorporationParams{
		WebRefreshToken: staticdb.WEB_REFRESH_TOKEN,
		CorporationId:   staticdb.CORPORATION_ID,
	})
	if err != nil {
		return nil, err
	}

	// // filter for buyback and shop contracts and insert them
	contracts := newContracts()
	for i := 0; i < hrwc.NumPages; i++ {

		// receive the next page
		page, err := hrwc.RecvUpdateExpires()
		if err != nil {
			return nil, err
		}

		for _, entry := range page {
			// append the order to the type orders
			contracts.filterAddEntry(
				staticdb.CORPORATION_ID,
				entry,
			)
		}
	} // //

	return cache.NewExpirableDataPtr(
		*contracts,
		hrwc.Expires,
	), nil
}
