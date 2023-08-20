package db

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/authing"
)

type AppraisalParams string

func (ap AppraisalParams) CacheKey() string {
	return fmt.Sprintf("appraisal-%s", ap)
}

type WC_ShopAppraisalClient = client.CachingClient[
	AppraisalParams,
	ShopAppraisal,
	cache.ExpirableData[ShopAppraisal],
	ShopAppraisalClient,
]

type WC_BuybackAppraisalClient = client.CachingClient[
	AppraisalParams,
	BuybackAppraisal,
	cache.ExpirableData[BuybackAppraisal],
	BuybackAppraisalClient,
]

type ShopAppraisalClient struct {
	*DatabaseClient
	expires time.Duration
}

func (sac ShopAppraisalClient) Fetch(
	ctx context.Context,
	params AppraisalParams,
) (*cache.ExpirableData[ShopAppraisal], error) {
	if appraisal, err := sac.GetShopAppraisal(
		ctx,
		string(params),
	); err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr[ShopAppraisal](
			appraisal,
			time.Now().Add(sac.expires),
		), nil
	}
}

type BuybackAppraisalClient struct {
	*DatabaseClient
	expires time.Duration
}

func (bac BuybackAppraisalClient) Fetch(
	ctx context.Context,
	params AppraisalParams,
) (*cache.ExpirableData[BuybackAppraisal], error) {
	if appraisal, err := bac.GetBuybackAppraisal(
		ctx,
		string(params),
	); err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr[BuybackAppraisal](
			appraisal,
			time.Now().Add(bac.expires),
		), nil
	}
}

type ShopQueueParams struct {
	CharacterRefreshToken string
}

func (sqp ShopQueueParams) AuthRefreshToken() string {
	return sqp.CharacterRefreshToken
}

func (ShopQueueParams) CacheKey() string {
	return "shopqueue"
}

type A_SC_ShopQueueClient = authing.AuthingClient[
	ShopQueueParams,
	client.CachingRep[[]string],
	client.StrongCachingClient[
		ShopQueueParams,
		[]string,
		cache.ExpirableData[[]string],
		ShopQueueClient,
	],
]

type ShopQueueClient struct {
	*DatabaseClient
	expires time.Duration
}

func (sqrc ShopQueueClient) Fetch(
	ctx context.Context,
	params ShopQueueParams,
) (*cache.ExpirableData[[]string], error) {
	if queue, err := sqrc.GetShopQueue(ctx); err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr[[]string](
			queue,
			time.Now().Add(sqrc.expires),
		), nil
	}
}

type CharacterCodesParams struct {
	CharacterId int32
}

func (ccp CharacterCodesParams) CacheKey() string {
	return fmt.Sprintf("charcodes-%d", ccp.CharacterId)
}

type CachingCharacterCodesClient = client.StrongCachingClient[
	CharacterCodesParams,
	CharacterCodes,
	cache.ExpirableData[CharacterCodes],
	CharacterCodesClient,
]

type CharacterCodesClient struct {
	*DatabaseClient
	expires time.Duration
}

func (ccrc CharacterCodesClient) Fetch(
	ctx context.Context,
	params CharacterCodesParams,
) (*cache.ExpirableData[CharacterCodes], error) {
	if codesRep, err := ccrc.GetCharacterCodes(
		ctx,
		params.CharacterId,
	); err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr[CharacterCodes](
			*codesRep,
			time.Now().Add(ccrc.expires),
		), nil
	}
}
