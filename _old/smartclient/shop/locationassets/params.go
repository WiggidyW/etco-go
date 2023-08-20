package locationassets

import "fmt"

type ShopLocationAssetsParams struct {
	CorporationId int32
	RefreshToken  string
	LocationId    int64
}

func (p ShopLocationAssetsParams) CacheKey() string {
	return fmt.Sprintf(
		"shopassets-%d-%d",
		p.CorporationId,
		p.LocationId,
	)
}
