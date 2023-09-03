package allassets_

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	massetscorporation "github.com/WiggidyW/etco-go/client/esi/model/assetscorporation"
)

const (
	ALL_SHOP_ASSETS_MIN_EXPIRES    time.Duration = 0
	ALL_SHOP_ASSETS_SLOCK_TTL      time.Duration = 30 * time.Second
	ALL_SHOP_ASSETS_SLOCK_MAX_WAIT time.Duration = 15 * time.Second
)

type WC_AllShopAssetsClient = wc.WeakCachingClient[
	AllShopAssetsParams,
	map[int64]map[int32]*int64,
	cache.ExpirableData[map[int64]map[int32]*int64],
	AllShopAssetsClient,
]

func NewWC_AllShopAssetsClient(
	modelacClient massetscorporation.AssetsCorporationClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_AllShopAssetsClient {
	return wc.NewWeakCachingClient(
		NewAllShopAssetsClient(modelacClient),
		ALL_SHOP_ASSETS_MIN_EXPIRES,
		cCache,
		sCache,
		ALL_SHOP_ASSETS_SLOCK_TTL,
		ALL_SHOP_ASSETS_SLOCK_MAX_WAIT,
	)
}

type AllShopAssetsClient struct {
	modelacClient massetscorporation.AssetsCorporationClient
}

func NewAllShopAssetsClient(
	modelacClient massetscorporation.AssetsCorporationClient,
) AllShopAssetsClient {
	return AllShopAssetsClient{modelacClient}
}

// TODO: add multi-caching supportclient for multiple locations
func (sac AllShopAssetsClient) Fetch(
	ctx context.Context,
	params AllShopAssetsParams,
) (*cache.ExpirableData[map[int64]map[int32]*int64], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := sac.modelacClient.Fetch(
		ctx,
		massetscorporation.AssetsCorporationParams{
			WebRefreshToken: build.CORPORATION_WEB_REFRESH_TOKEN,
			CorporationId:   build.CORPORATION_ID,
		},
	)
	if err != nil {
		return nil, err
	}

	// // insert all entries into unflattened assets
	unflattenedAssets := newUnflattenedAssets(hrwc.NumPages *
		massetscorporation.ASSETS_CORPORATION_ENTRIES_PER_PAGE)
	for i := 0; i < hrwc.NumPages; i++ {

		// receive the next page
		page, err := hrwc.RecvUpdateExpires()
		if err != nil {
			return nil, err
		}

		// add the entries to unflattened assets
		for _, entry := range page {
			unflattenedAssets.addEntry(entry)
		}
	} // //

	// flatten and filter the entries, converting them to shop assets
	shopAssets := unflattenedAssets.flattenAndFilter()

	return cache.NewExpirableDataPtr(
		shopAssets,
		hrwc.Expires,
	), nil
}
