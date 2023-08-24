package location

import (
	"context"

	// "github.com/WiggidyW/eve-trading-co-go/cache"
	// sc "github.com/WiggidyW/eve-trading-co-go/client/caching/strong/caching"
	// "github.com/WiggidyW/eve-trading-co-go/client/inventory/data"
	"github.com/WiggidyW/eve-trading-co-go/client/inventory/internal/unreserved"
)

// type SC_LocationShopAssetsClient = sc.StrongCachingClient[
// 	LocationShopAssetsParams,
// 	map[int32]inventory.ShopAsset,
// 	cache.ExpirableData[map[int32]inventory.ShopAsset],
// 	LocationShopAssetsClient,
// ]

type LocationShopAssetsClient struct {
	Inner unreserved.SC_UnreservedShopAssetsClient
}

func (lsac LocationShopAssetsClient) Fetch(
	ctx context.Context,
	params LocationShopAssetsParams,
	// ) (*cache.ExpirableData[map[int32]inventory.ShopAsset], error) {
) (map[int32]int64, error) {
	if unresRep, err := lsac.Inner.Fetch(
		ctx,
		unreserved.UnreservedShopAssetsParams{
			ShopQueue: params.ShopQueue,
		},
	); err != nil {
		return nil, err
	} else {
		return unresRep.Data()[params.LocationId], nil
		// return cache.NewExpirableDataPtr(
		// 	unresRep.Data()[params.LocationId],
		// 	unresRep.Expires(),
		// ), nil
	}
}
