package locationassets

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	massetscorporation "github.com/WiggidyW/etco-go/client/esi/model/assetscorporation"
	"github.com/WiggidyW/etco-go/client/inventory/locationassets/unreservedassets_"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
)

type LocationShopAssetsClient struct {
	unreservedClient unreservedassets_.SC_UnreservedShopAssetsClient
}

func NewLocationShopAssetsClient(
	modelacClient massetscorporation.AssetsCorporationClient,
	appraisalClient rdbc.WC_ReadShopAppraisalClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) LocationShopAssetsClient {
	return LocationShopAssetsClient{
		unreservedassets_.NewSC_UnreservedShopAssetsClient(
			modelacClient,
			appraisalClient,
			cCache,
			sCache,
		),
	}
}

func (
	lsac LocationShopAssetsClient,
) GetUnreservedAntiCache() *cache.StrongAntiCache {
	return lsac.unreservedClient.GetAntiCache()
}

func (lsac LocationShopAssetsClient) Fetch(
	ctx context.Context,
	params LocationShopAssetsParams,
) (map[int32]int64, error) {
	if unresRep, err := lsac.unreservedClient.Fetch(
		ctx,
		unreservedassets_.UnreservedShopAssetsParams{
			ShopQueue: params.ShopQueue,
		},
	); err != nil {
		return nil, err
	} else {
		return unresRep.Data()[params.LocationId], nil
	}
}
