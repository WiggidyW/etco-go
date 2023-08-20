package shopitems

import (
	"github.com/WiggidyW/weve-esi/client/smartclient/market"
	"github.com/WiggidyW/weve-esi/client/smartclient/shop/assets"
	"github.com/WiggidyW/weve-esi/staticdb"
)

type ShopItem struct {
	Naming staticdb.Naming
	Asset  assets.ShopAsset
	Price  market.ShopPrice
}
