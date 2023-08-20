package single

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	i "github.com/WiggidyW/weve-esi/client/contracts/items"
	cimodel "github.com/WiggidyW/weve-esi/client/esi/model/contractitems"
	"github.com/WiggidyW/weve-esi/staticdb"
)

type WC_RateLimitingContractItemsClient = wc.WeakCachingClient[
	RateLimitingContractItemsParams,
	[]i.ContractItem,
	cache.ExpirableData[[]i.ContractItem],
	RateLimitingContractItemsClient,
]

// rate limiting client
// https://github.com/esi/esi-issues/issues/636
type RateLimitingContractItemsClient struct {
	Inner       cimodel.ContractItemsClient
	MaxAttempts int
}

func (rlcic RateLimitingContractItemsClient) Fetch(
	ctx context.Context,
	params RateLimitingContractItemsParams,
) (*cache.ExpirableData[[]i.ContractItem], error) {
	modelRep, err := rlcic.fetchAttempt(
		ctx,
		cimodel.ContractItemsParams{
			WebRefreshToken: staticdb.WEB_REFRESH_TOKEN,
			CorporationId:   staticdb.CORPORATION_ID,
			ContractId:      params.ContractId,
		},
		1,
	)
	if err != nil {
		return nil, err
	}

	return cache.NewExpirableDataPtr(
		i.EntriesToItems(modelRep.Data()),
		modelRep.Expires(),
	), nil
}

func (rlcic RateLimitingContractItemsClient) fetchAttempt(
	ctx context.Context,
	params cimodel.ContractItemsParams,
	attempt int,
) (*cache.ExpirableData[[]cimodel.ContractItemsEntry], error) {
	modelRep, err := rlcic.Inner.Fetch(ctx, params)
	if err != nil {
		if attempt < rlcic.MaxAttempts && i.RateLimited(err) {
			attempt++
			time.Sleep(i.LIMITED_SLEEP)
			return rlcic.fetchAttempt(ctx, params, attempt+1)
		} else { // out of attempts or not rate limited
			return nil, err
		}
	}

	return modelRep, nil
}
