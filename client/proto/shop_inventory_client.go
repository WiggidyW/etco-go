package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"
	"github.com/WiggidyW/etco-go/client/inventory"
	"github.com/WiggidyW/etco-go/client/market"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBShopInventoryParams struct {
	TypeNamingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap]
	LocationId        int64
}

type PBShopInventoryClient struct {
	rInventoryClient inventory.InventoryClient
	rShopPriceClient market.ShopPriceClient
}

func (sic PBShopInventoryClient) Fetch(
	ctx context.Context,
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
	rInventory, err := sic.fetchRInventory(ctx, params.LocationId)
	if err != nil {
		return nil, err
	}

	// fetch a shop item for each item in the inventory
	chnSendShopItem, chnRecvShopItem := chanresult.
		NewChanResult[*proto.ShopItem](ctx, len(rInventory), 0).Split()
	for typeId, quantity := range rInventory {
		go sic.transceiveFetchShopItem(
			ctx,
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
	ctx context.Context,
	locationId int64,
) (
	rInventory map[int32]int64,
	err error,
) {
	rInventoryPtr, err := sic.rInventoryClient.Fetch(
		ctx,
		inventory.InventoryParams{LocationId: locationId},
	)
	if err != nil {
		return nil, err
	} else {
		return *rInventoryPtr, nil
	}
}

func (sic PBShopInventoryClient) transceiveFetchShopItem(
	ctx context.Context,
	typeId int32,
	quantity int64,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	shopLocationInfo staticdb.ShopLocationInfo,
	chnSend chanresult.ChanSendResult[*proto.ShopItem],
) error {
	shopItem, err := sic.fetchShopItem(
		ctx,
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
	ctx context.Context,
	typeId int32,
	quantity int64,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	shopLocationInfo staticdb.ShopLocationInfo,
) (
	item *proto.ShopItem,
	err error,
) {
	rShopItem, err := sic.rShopPriceClient.Fetch(
		ctx,
		market.ShopPriceParams{
			ShopLocationInfo: shopLocationInfo,
			TypeId:           typeId,
			Quantity:         quantity,
		},
	)
	if err != nil {
		return nil, err
	} else {
		return pu.NewPBShopItem(
			*rShopItem,
			namingSession,
		), nil
	}
}
