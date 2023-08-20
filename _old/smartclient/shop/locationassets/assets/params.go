package assets

import "fmt"

type ShopAssetsParams struct {
	CorporationId int32
	RefreshToken  string
}

func (f ShopAssetsParams) CacheKey() string {
	return fmt.Sprintf("shopassets-%d", f.CorporationId)
}
