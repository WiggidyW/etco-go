package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
)

const (
	CAPACITY_MULTIPLIER = 3
	CAPACITY_DIVISOR    = 2
)

type BucketClient struct{ *b.BucketClient }

func NewBucketClient(nameSpace string, creds []byte) BucketClient {
	return BucketClient{b.NewBucketClient(nameSpace, creds)}
}

func (bc BucketClient) ReadWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
) (map[b.TypeId]b.WebBuybackSystemTypeBundle, error) {
	return bc.BucketClient.ReadWebBuybackSystemTypeMapsBuilder(
		ctx,
		transformWebCapacity(
			build.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
		),
	)
}

func (bc BucketClient) ReadWebShopLocationTypeMapsBuilder(
	ctx context.Context,
) (map[b.TypeId]b.WebShopLocationTypeBundle, error) {
	return bc.BucketClient.ReadWebShopLocationTypeMapsBuilder(
		ctx,
		transformWebCapacity(
			build.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
		),
	)
}

func (bc BucketClient) ReadWebBuybackSystems(
	ctx context.Context,
) (map[b.SystemId]b.WebBuybackSystem, error) {
	return bc.BucketClient.ReadWebBuybackSystems(
		ctx,
		transformWebCapacity(build.CAPACITY_WEB_BUYBACK_SYSTEMS),
	)
}

func (bc BucketClient) ReadWebShopLocations(
	ctx context.Context,
) (map[b.LocationId]b.WebShopLocation, error) {
	return bc.BucketClient.ReadWebShopLocations(
		ctx,
		transformWebCapacity(build.CAPACITY_WEB_SHOP_LOCATIONS),
	)
}

func (bc BucketClient) ReadWebMarkets(
	ctx context.Context,
) (map[b.MarketName]b.WebMarket, error) {
	return bc.BucketClient.ReadWebMarkets(
		ctx,
		transformWebCapacity(build.CAPACITY_WEB_MARKETS),
	)
}

func (bc BucketClient) ReadAuthHashSet(
	ctx context.Context,
	key string,
) (b.AuthHashSet, error) {
	return bc.BucketClient.ReadAuthHashSet(ctx, key)
}

func (bc BucketClient) WriteWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	v map[b.TypeId]b.WebBuybackSystemTypeBundle,
) error {
	return bc.BucketClient.WriteWebBuybackSystemTypeMapsBuilder(ctx, v)
}

func (bc BucketClient) WriteWebShopLocationTypeMapsBuilder(
	ctx context.Context,
	v map[b.TypeId]b.WebShopLocationTypeBundle,
) error {
	return bc.BucketClient.WriteWebShopLocationTypeMapsBuilder(ctx, v)
}

func (bc BucketClient) WriteWebShopLocations(
	ctx context.Context,
	v map[b.LocationId]b.WebShopLocation,
) error {
	return bc.BucketClient.WriteWebShopLocations(ctx, v)
}

func (bc BucketClient) WriteWebBuybackSystems(
	ctx context.Context,
	v map[b.SystemId]b.WebBuybackSystem,
) error {
	return bc.BucketClient.WriteWebBuybackSystems(ctx, v)
}

func (bc BucketClient) WriteWebMarkets(
	ctx context.Context,
	v map[b.MarketName]b.WebMarket,
) error {
	return bc.BucketClient.WriteWebMarkets(ctx, v)
}

func (bc BucketClient) WriteAuthHashSet(
	ctx context.Context,
	v b.AuthHashSet,
	key string,
) error {
	return bc.BucketClient.WriteAuthHashSet(ctx, v, key)
}
