package assets

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
)

type CachingShopLocationAssetsClient = *client.CachingClient[
	ShopLocationAssetsParams,
	[]ShopAsset,
	cache.ExpirableData[[]ShopAsset],
	ShopLocationAssetsClient,
]

type ShopLocationAssetsParams struct {
	CorporationId int32
	RefreshToken  string
	LocationId    int64
}

func (p ShopLocationAssetsParams) CacheKey() string {
	return fmt.Sprintf(
		"shopassets-%d-%d",
		p.CorporationId,
		p.LocationId,
	)
}

type ShopLocationAssetsClient struct {
	client CachingShopAssetsClient
}

// fetches the map[location][]asset and returns the []asset for the provided location
// if it's missing, we return it anyways as a nil slice
func (slac ShopLocationAssetsClient) Fetch(
	ctx context.Context,
	params ShopLocationAssetsParams,
) (*cache.ExpirableData[[]ShopAsset], error) {
	if assetsRep, err := slac.client.Fetch(ctx, ShopAssetsParams{
		CorporationId: params.CorporationId,
		RefreshToken:  params.RefreshToken,
	}); err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			assetsRep.Data()[params.LocationId],
			assetsRep.Expires(),
		), nil
	}
}
