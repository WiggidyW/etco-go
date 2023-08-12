package assets

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/modelclient"
)

type CachingShopAssetsClient = *client.CachingClient[
	ShopAssetsParams,
	map[int64][]ShopAsset,
	cache.ExpirableData[map[int64][]ShopAsset],
	ShopAssetsClient,
]

type ShopAssetsParams struct {
	CorporationId int32
	RefreshToken  string
}

func (f ShopAssetsParams) CacheKey() string {
	return fmt.Sprintf("shopassets-%d", f.CorporationId)
}

type ShopAssetsClient struct {
	client *modelclient.ClientAssetsCorporation
}

// TODO: add multi-caching supportclient for multiple locations
func (sac ShopAssetsClient) Fetch(
	ctx context.Context,
	params ShopAssetsParams,
) (*cache.ExpirableData[map[int64][]ShopAsset], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the ESI assets
	chnRecv, err := sac.client.FetchStreamBlocking(
		ctx,
		modelclient.NewFetchParamsAssetsCorporation(
			params.CorporationId,
			params.RefreshToken,
		),
	)
	if err != nil {
		return nil, err
	}

	// initialize unflattened assets
	unflattenedAssets := newUnflattenedAssets(
		chnRecv.NumPages() * sac.client.EntriesPerPage(),
	)

	// handle the pages
	var maxExpires time.Time = chnRecv.HeadExpires()
	for numPages := chnRecv.NumPages(); numPages > 0; numPages-- {
		page, err := chnRecv.Recv()
		if err != nil {
			return nil, err
		}
		// update the maxExpires
		if page.Expires().After(maxExpires) {
			maxExpires = page.Expires()
		}
		// add the entries to unflattened assets
		for _, entry := range page.Data() {
			unflattenedAssets.addEntry(entry)
		}
	}

	// flatten and filter the entries, converting them to shop assets
	shopAssets := unflattenedAssets.flattenAndFilter()

	// return the shop assets and their expiration date
	data := cache.NewExpirableData(shopAssets, maxExpires)
	return &data, nil
}
