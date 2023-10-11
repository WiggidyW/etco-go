package contracts

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	mcontractscorporation "github.com/WiggidyW/etco-go/client/esi/model/contractscorporation"
)

const (
	CONTRACTS_MIN_EXPIRES   time.Duration = 0
	CONTRACTS_SLOCK_TTL     time.Duration = 1 * time.Minute
	CONTRACTS_SLOCK_MAXWAIT time.Duration = 30 * time.Second
)

type WC_ContractsClient = wc.WeakCachingClient[
	ContractsParams,
	Contracts,
	cache.ExpirableData[Contracts],
	ContractsClient,
]

func NewWC_ContractsClient(
	modelClient mcontractscorporation.ContractsCorporationClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_ContractsClient {
	return wc.NewWeakCachingClient(
		NewContractsClient(modelClient),
		CONTRACTS_MIN_EXPIRES,
		cCache,
		sCache,
		CONTRACTS_SLOCK_TTL,
		CONTRACTS_SLOCK_MAXWAIT,
	)
}

type ContractsClient struct {
	modelClient mcontractscorporation.ContractsCorporationClient
}

func NewContractsClient(
	modelClient mcontractscorporation.ContractsCorporationClient,
) ContractsClient {
	return ContractsClient{modelClient}
}

func (cc ContractsClient) Fetch(
	ctx context.Context,
	params ContractsParams,
) (*cache.ExpirableData[Contracts], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := cc.modelClient.Fetch(
		ctx,
		mcontractscorporation.ContractsCorporationParams{
			WebRefreshToken: build.CORPORATION_WEB_REFRESH_TOKEN,
			CorporationId:   build.CORPORATION_ID,
		},
	)
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
				build.CORPORATION_ID,
				entry,
			)
		}
	} // //

	return cache.NewExpirableDataPtr(
		*contracts,
		hrwc.Expires,
	), nil
}
