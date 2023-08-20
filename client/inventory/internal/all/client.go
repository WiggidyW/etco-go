package all

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	m "github.com/WiggidyW/weve-esi/client/esi/model/assetscorporation"
	"github.com/WiggidyW/weve-esi/staticdb"
)

type WC_AllShopAssetsClient = *wc.WeakCachingClient[
	AllShopAssetsParams,
	map[int64]map[int32]*int64,
	cache.ExpirableData[map[int64]map[int32]*int64],
	AllShopAssetsClient,
]

type AllShopAssetsClient struct {
	Inner m.AssetsCorporationClient
}

// TODO: add multi-caching supportclient for multiple locations
func (sac AllShopAssetsClient) Fetch(
	ctx context.Context,
	params AllShopAssetsParams,
) (*cache.ExpirableData[map[int64]map[int32]*int64], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := sac.Inner.Fetch(ctx, m.AssetsCorporationParams{
		WebRefreshToken: staticdb.WEB_REFRESH_TOKEN,
		CorporationId:   staticdb.CORPORATION_ID,
	})
	if err != nil {
		return nil, err
	}

	// // insert all entries into unflattened assets
	unflattenedAssets := newUnflattenedAssets(
		hrwc.NumPages * m.ASSETS_CORPORATION_ENTRIES_PER_PAGE,
	)
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
