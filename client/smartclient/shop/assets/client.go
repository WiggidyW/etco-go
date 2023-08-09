package assets

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client/modelclient"
	"github.com/WiggidyW/weve-esi/staticdb/tc"
)

type ShopAssetsClientFetchParams struct {
	CorporationId int32
	RefreshToken  string
	ShopInfo      *tc.ShopInfo
}

func (f ShopAssetsClientFetchParams) Key() string {
	return fmt.Sprintf("shopassets-%d", f.CorporationId)
}

type ShopAssetsClient struct {
	client *modelclient.ClientAssetsCorporation
}

func (fac *ShopAssetsClient) Fetch(
	ctx context.Context,
	params ShopAssetsClientFetchParams,
) (*cache.ExpirableData[map[int64][]ShopAsset], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the ESI assets
	strm := fac.client.FetchStreamBlocking(
		ctx,
		modelclient.NewFetchParamsAssetsCorporation(
			params.CorporationId,
			params.RefreshToken,
		),
	)
	defer strm.Close()

	// initialize unflattened assets
	unflattenedAssets := newUnflattenedAssets(
		strm.NumPages()*fac.client.EntriesPerPage(),
		params.ShopInfo,
	)

	// handle the pages
	var maxExpires time.Time = time.Unix(0, 0)
	for pages := strm.NumPages(); pages > 0; pages-- {
		page, err := strm.Recv()
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
