package shopassets

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/staticdb"
)

const (
	RAW_BUF_CAP        int = 0
	UNRESERVED_BUF_CAP int = 0
)

func init() {
	keys.TypeStrNSRawShopAssets = cache.RegisterType[struct{}]("rawshopassets", 0)
	keys.TypeStrRawShopAssets = cache.RegisterType[map[int32]int64]("shopassets", RAW_BUF_CAP)
	keys.TypeStrUnreservedShopAssets = cache.RegisterType[map[int32]int64]("unreservedshopassets", UNRESERVED_BUF_CAP)
}

func getRawShopAssets(
	x cache.Context,
	locationId int64,
) (
	rep map[int32]int64,
	expires time.Time,
	err error,
) {
	return rawShopAssetsGet(x, locationId)
}

func GetUnreservedShopAssets(
	x cache.Context,
	locationId int64,
) (
	assets map[int32]int64,
	expires time.Time,
	err error,
) {
	return unreservedShopAssetsGet(x, locationId)
}

type ProtoShopInventoryRep struct {
	LocationInfo *proto.LocationInfo
	Assets       []*proto.ShopItem
}

func ProtoGetShopInventory(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	locationId int64,
) (
	rep ProtoShopInventoryRep,
	expires time.Time,
	err error,
) {
	shopLocationInfo := staticdb.GetShopLocationInfo(locationId)
	if shopLocationInfo == nil {
		err = protoerr.MsgNew(protoerr.NOT_FOUND, "Shop Location not valid")
		return rep, expires, err
	}

	// fetch unreserved shop assets
	var rAssets map[int32]int64
	rAssets, expires, err = GetUnreservedShopAssets(x, locationId)
	if err != nil {
		return rep, expires, err
	}

	// fetch shop item for each asset in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chn := expirable.NewChanResult[*proto.ShopItem](x.Ctx(), len(rAssets), 0)
	for typeId, quantity := range rAssets {
		go expirable.P5Transceive(
			chn,
			x, r, typeId, quantity, *shopLocationInfo,
			market.ProtoGetShopPrice,
		)
	}

	// fetch proto location info
	rep.LocationInfo, expires, err =
		esi.ProtoGetLocationInfoCOV(x, r, locationId).RecvExpMin(expires)
	if err != nil {
		return rep, expires, err
	}

	// recv shop items
	rep.Assets = make([]*proto.ShopItem, 0, len(rAssets))
	for i := 0; i < len(rAssets); i++ {
		var asset *proto.ShopItem
		asset, expires, err = chn.RecvExpMin(expires)
		if err != nil {
			rep.Assets = nil
			return rep, expires, err
		} else if asset != nil && asset.PricePerUnit > 0.0 {
			rep.Assets = append(rep.Assets, asset)
		}
	}

	return rep, expires, nil
}
