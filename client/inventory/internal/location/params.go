package location

type LocationShopAssetsParams struct {
	ShopQueue  []string
	LocationId int64
}

func NewLocationShopAssetsParams(
	shopQueue []string,
	locationId int64,
) LocationShopAssetsParams {
	return LocationShopAssetsParams{
		ShopQueue:  shopQueue,
		LocationId: locationId,
	}
}
