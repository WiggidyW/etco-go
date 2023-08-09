package shopitems

import (
	"github.com/WiggidyW/weve-esi/client/smartclient/market"
	"github.com/WiggidyW/weve-esi/client/smartclient/shop/assets"
	"github.com/WiggidyW/weve-esi/staticdb/sde"
)

type ShopItem struct {
	Naming sde.Naming
	Asset  assets.ShopAsset
	Price  market.ShopPrice
}
