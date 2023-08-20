package locationassets

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/shop/locationassets/assets"
)

type CachingShopLocationAssetsClient = *client.CachingClient[
	ShopLocationAssetsParams,
	[]ShopAsset,
	cache.ExpirableData[[]ShopAsset],
	ShopLocationAssetsClient,
]

type ShopLocationAssetsClient struct {
	client assets.ShopAssetsClient
}

// fetches the map[location][]asset and returns the []asset for the provided location
// if it's missing, we return it anyways as a nil slice
func (slac ShopLocationAssetsClient) Fetch(
	ctx context.Context,
	params ShopLocationAssetsParams,
) (*cache.ExpirableData[[]ShopAsset], error) {
	if assetsRep, err := slac.client.Fetch(ctx, assets.ShopAssetsParams{
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
