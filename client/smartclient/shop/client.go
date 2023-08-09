package shopitems

import (
	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	shopassets "github.com/WiggidyW/weve-esi/client/smartclient/shop/assets"
)

type ShopItemsClientFetchParams struct {
	LocationId    int64
	CorporationId int32
	RefreshToken  string
}

type ShopItemsClient struct {
	assetClient *client.CachingClient[
		shopassets.ShopAssetsClientFetchParams,
		map[int64][]shopassets.ShopAsset,
		cache.ExpirableData[map[int64][]shopassets.ShopAsset],
		*shopassets.ShopAssetsClient,
	]
}
