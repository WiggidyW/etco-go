package proto

import (
	"github.com/WiggidyW/chanresult"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/shopassets"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBShopInventoryParams struct {
	TypeNamingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap]
	LocationId        int64
}

type PBShopInventoryClient struct{}

func NewPBShopInventoryClient() PBShopInventoryClient {
	return PBShopInventoryClient{}
}

func (sic PBShopInventoryClient) Fetch(
	x cache.Context,
	params PBShopInventoryParams,
) (
	items []*proto.ShopItem,
	err error,
) {
	// if the location has no shop, return empty inventory
	shopLocationInfoPtr := staticdb.GetShopLocationInfo(params.LocationId)
	if shopLocationInfoPtr == nil {
		return items, nil
	}
	shopLocationInfo := *shopLocationInfoPtr

	// fetch the raw inventory
	rInventory, err := sic.fetchRInventory(x, params.LocationId)
	if err != nil {
		return nil, err
	}

	// fetch a shop item for each item in the inventory
	chnSendShopItem, chnRecvShopItem := chanresult.
		NewChanResult[*proto.ShopItem](x.Ctx(), len(rInventory), 0).Split()
	for typeId, quantity := range rInventory {
		go sic.transceiveFetchShopItem(
			x,
			typeId,
			quantity,
			params.TypeNamingSession,
			shopLocationInfo,
			chnSendShopItem,
		)
	}

	// collect all shop items
	items = make([]*proto.ShopItem, 0, len(rInventory))
	for i := 0; i < len(rInventory); i++ {
		shopItem, err := chnRecvShopItem.Recv()
		if err != nil {
			return nil, err
		} else {
			items = append(items, shopItem)
		}
	}

	return items, nil
}

func (sic PBShopInventoryClient) fetchRInventory(
	x cache.Context,
	locationId int64,
) (
	rInventory map[int32]int64,
	err error,
) {
	rInventory, _, err = shopassets.GetUnreservedShopAssets(x, locationId)
	return rInventory, err
}

func (sic PBShopInventoryClient) transceiveFetchShopItem(
	x cache.Context,
	typeId int32,
	quantity int64,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	shopLocationInfo staticdb.ShopLocationInfo,
	chnSend chanresult.ChanSendResult[*proto.ShopItem],
) error {
	shopItem, err := sic.fetchShopItem(
		x,
		typeId,
		quantity,
		namingSession,
		shopLocationInfo,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(shopItem)
	}
}

func (sic PBShopInventoryClient) fetchShopItem(
	x cache.Context,
	typeId int32,
	quantity int64,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	shopLocationInfo staticdb.ShopLocationInfo,
) (
	item *proto.ShopItem,
	err error,
) {
	rShopItem, _, err := market.GetShopPrice(
		x,
		typeId,
		quantity,
		shopLocationInfo,
	)
	if err != nil {
		return nil, err
	} else {
		return pu.NewPBShopItem(
			rShopItem,
			namingSession,
		), nil
	}
}
