package etcogobucket

import (
	"context"
)

func (bc *BucketClient) ReadAttrsWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
) (*Attrs, error) {
	return bc.readAttrs(
		ctx,
		WEB,
		OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}

func (bc *BucketClient) ReadAttrsWebShopLocationTypeMapsBuilder(
	ctx context.Context,
) (*Attrs, error) {
	return bc.readAttrs(
		ctx,
		WEB,
		OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}

func (bc *BucketClient) ReadAttrsWebBuybackSystems(
	ctx context.Context,
) (*Attrs, error) {
	return bc.readAttrs(
		ctx,
		WEB,
		OBJNAME_WEB_BUYBACK_SYSTEMS,
	)
}

func (bc *BucketClient) ReadAttrsWebShopLocations(
	ctx context.Context,
) (*Attrs, error) {
	return bc.readAttrs(
		ctx,
		WEB,
		OBJNAME_WEB_SHOP_LOCATIONS,
	)
}

func (bc *BucketClient) ReadAttrsWebMarkets(
	ctx context.Context,
) (*Attrs, error) {
	return bc.readAttrs(
		ctx,
		WEB,
		OBJNAME_WEB_MARKETS,
	)
}

func (bc *BucketClient) ReadWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	capacity int,
) (v map[TypeId]WebBuybackSystemTypeBundle, err error) {
	v = make(
		map[TypeId]WebBuybackSystemTypeBundle,
		capacity,
	)
	_, err = read(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
		&v,
	)
	return v, err
}

func (bc *BucketClient) ReadWebShopLocationTypeMapsBuilder(
	ctx context.Context,
	capacity int,
) (v map[TypeId]WebShopLocationTypeBundle, err error) {
	v = make(
		map[TypeId]WebShopLocationTypeBundle,
		capacity,
	)
	_, err = read(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
		&v,
	)
	return v, err
}

func (bc *BucketClient) ReadWebBuybackSystems(
	ctx context.Context,
	capacity int,
) (v map[SystemId]WebBuybackSystem, err error) {
	v = make(
		map[SystemId]WebBuybackSystem,
		capacity,
	)
	_, err = read(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_BUYBACK_SYSTEMS,
		&v,
	)
	return v, err
}

func (bc *BucketClient) ReadWebShopLocations(
	ctx context.Context,
	capacity int,
) (v map[LocationId]WebShopLocation, err error) {
	v = make(
		map[LocationId]WebShopLocation,
		capacity,
	)
	_, err = read(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_SHOP_LOCATIONS,
		&v,
	)
	return v, err
}

func (bc *BucketClient) ReadWebMarkets(
	ctx context.Context,
	capacity int,
) (v map[MarketName]WebMarket, err error) {
	v = make(
		map[MarketName]WebMarket,
		capacity,
	)
	_, err = read(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_MARKETS,
		&v,
	)
	return v, err
}

func (bc *BucketClient) WriteWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	v map[TypeId]WebBuybackSystemTypeBundle,
) error {
	return write(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
		v,
	)
}

func (bc *BucketClient) WriteWebShopLocationTypeMapsBuilder(
	ctx context.Context,
	v map[TypeId]WebShopLocationTypeBundle,
) error {
	return write(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
		v,
	)
}

func (bc *BucketClient) WriteWebShopLocations(
	ctx context.Context,
	v map[LocationId]WebShopLocation,
) error {
	return write(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_SHOP_LOCATIONS,
		v,
	)
}

func (bc *BucketClient) WriteWebBuybackSystems(
	ctx context.Context,
	v map[SystemId]WebBuybackSystem,
) error {
	return write(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_BUYBACK_SYSTEMS,
		v,
	)
}

func (bc *BucketClient) WriteWebMarkets(
	ctx context.Context,
	v map[MarketName]WebMarket,
) error {
	return write(
		bc,
		ctx,
		WEB,
		OBJNAME_WEB_MARKETS,
		v,
	)
}
