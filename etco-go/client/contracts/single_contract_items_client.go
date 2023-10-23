package contracts

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	cimodel "github.com/WiggidyW/etco-go/client/esi/model/contractitems"
)

const (
	CONTRACT_ITEMS_MAX_ATTEMPTS  int           = 3
	CONTRACT_ITEMS_MIN_EXPIRES   time.Duration = 48 * time.Hour
	CONTRACT_ITEMS_SLOCK_TTL     time.Duration = 30 * time.Second
	CONTRACT_ITEMS_SLOCK_MAXWAIT time.Duration = 10 * time.Second
)

type WC_SingleContractItemsClient = wc.WeakCachingClient[
	SingleContractItemsParams,
	[]ContractItem,
	cache.ExpirableData[[]ContractItem],
	*SingleContractItemsClient,
]

func NewWC_SingleContractItemsClient(
	modelClient cimodel.ContractItemsClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_SingleContractItemsClient {
	return wc.NewWeakCachingClient(
		NewSingleContractItemsClient(modelClient),
		CONTRACT_ITEMS_MIN_EXPIRES,
		cCache,
		sCache,
		CONTRACT_ITEMS_SLOCK_TTL,
		CONTRACT_ITEMS_SLOCK_MAXWAIT,
	)
}

// rate limiting client
// https://github.com/esi/esi-issues/issues/636
type SingleContractItemsClient struct {
	modelClient cimodel.ContractItemsClient
	maxAttempts int
	rateLimiter chan struct{} // capacity = REQS_PER_INTERVAL, fill in constructor
}

func NewSingleContractItemsClient(
	modelClient cimodel.ContractItemsClient,
) *SingleContractItemsClient {
	client := &SingleContractItemsClient{
		modelClient,
		CONTRACT_ITEMS_MAX_ATTEMPTS,
		make(chan struct{}, CI_REQS_PER_INTERVAL),
	}
	for i := 0; i < CI_REQS_PER_INTERVAL; i++ {
		client.rateLimiter <- struct{}{}
	}
	return client
}

func (scic *SingleContractItemsClient) Fetch(
	ctx context.Context,
	params SingleContractItemsParams,
) (*cache.ExpirableData[[]ContractItem], error) {
	modelRep, err := scic.fetchAttempt(
		ctx,
		cimodel.ContractItemsParams{
			WebRefreshToken: build.CORPORATION_WEB_REFRESH_TOKEN,
			CorporationId:   build.CORPORATION_ID,
			ContractId:      params.ContractId,
		},
		1,
	)
	if err != nil {
		return nil, err
	}

	return cache.NewExpirableDataPtr(
		EntriesToItems(modelRep.Data()),
		modelRep.Expires(),
	), nil
}

func (scic *SingleContractItemsClient) fetchAttempt(
	ctx context.Context,
	params cimodel.ContractItemsParams,
	attempt int,
) (*cache.ExpirableData[[]cimodel.ContractItemsEntry], error) {
	scic.rateLimiterStart()

	modelRep, err := scic.modelClient.Fetch(ctx, params)
	if err != nil && attempt < scic.maxAttempts && RateLimited(err) {
		scic.rateLimiterDone() // block
		return scic.fetchAttempt(ctx, params, attempt+1)
	}

	go scic.rateLimiterDone() // don't block

	if err != nil { // out of attempts or not rate limited
		return nil, err
	}

	return modelRep, nil
}

func (scic *SingleContractItemsClient) rateLimiterStart() {
	<-scic.rateLimiter
}

func (scic *SingleContractItemsClient) rateLimiterDone() {
	time.Sleep(CI_ATTEMPT_INTERVAL)
	scic.rateLimiter <- struct{}{}
}
