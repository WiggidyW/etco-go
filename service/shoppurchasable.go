package service

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/inventory"
	"github.com/WiggidyW/weve-esi/client/market/shop"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/util"
)

func (s *Service) ShopPurchasable(
	ctx context.Context,
	req *proto.ShopPurchasableRequest,
) (*proto.ShopPurchasableResponse, error) {
	inventoryRep, err := s.inventoryClient.Fetch(
		ctx,
		authingfwding.WithAuthableParams[inventory.InventoryParams]{
			NativeRefreshToken: req.Auth.Token,
			Params: inventory.InventoryParams{
				LocationId: req.LocationId,
			},
		},
	)

	ok, authRep, errRep := authRepToGrpcRep(inventoryRep, err)
	grpcRep := &proto.ShopPurchasableResponse{
		Auth:  authRep,
		Error: errRep,
	}
	if !ok {
		return grpcRep, nil
	}

	rInventory := *inventoryRep.Data
	syncNamingSession := maybeNewSyncNamingSession(req.IncludeNaming)

	shopLocationInfoPtr := staticdb.GetShopLocationInfo(req.LocationId)
	if shopLocationInfoPtr == nil {
		// this location has nothing to sell (no configuration)
		grpcRep.Items = []*proto.ShopItem{}
		grpcRep.Naming = maybeFinishNamingSession(syncNamingSession)
		return grpcRep, nil
	}
	shopLocationInfo := *shopLocationInfoPtr

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := util.NewChanResult[*proto.ShopItem](ctx).Split()

	// fetch all items (pricing + naming) in parallel
	for rTypeId, rQuantity := range rInventory {
		go s.fetchShopItem(
			ctx,
			rTypeId,
			rQuantity,
			shopLocationInfo,
			syncNamingSession,
			chnSend,
		)
	}

	// collect all items
	grpcRep.Items = make([]*proto.ShopItem, 0, len(rInventory))
	for i := 0; i < len(rInventory); i++ {
		if pbItem, err := chnRecv.Recv(); err != nil {
			grpcRep.Error = newErrorResponse(err)
			return grpcRep, nil
		} else if pbItem != nil {
			grpcRep.Items = append(grpcRep.Items, pbItem)
		}
	}

	grpcRep.Naming = maybeFinishNamingSession(syncNamingSession)

	return grpcRep, nil
}

func (s *Service) fetchShopItem(
	ctx context.Context,
	rTypeId int32,
	rQuantity int64,
	shopLocationInfo staticdb.ShopLocationInfo,
	syncNamingSession *staticdb.NamingSession[*staticdb.SyncIndexMap],
	chnSend util.ChanSendResult[*proto.ShopItem],
) error {
	if rPrice, err := s.shopPriceClient.Fetch(
		ctx,
		shop.ShopPriceParams{
			ShopLocationInfo: shopLocationInfo,
			TypeId:           rTypeId,
		},
	); err != nil {
		return chnSend.SendErr(err)
	} else if rPrice.PricePerUnit <= 0 {
		// send discard items with a rejected price
		return chnSend.SendOk(nil)
	} else {
		// send items with a valid price
		return chnSend.SendOk(newPBShopItem(
			rQuantity,
			*rPrice,
			syncNamingSession,
		))
	}
}

func newPBShopItem(
	rQuantity int64,
	rShopPrice shop.ShopPrice,
	syncNamingSession *staticdb.NamingSession[*staticdb.SyncIndexMap],
) *proto.ShopItem {
	return &proto.ShopItem{
		TypeId:       rShopPrice.TypeId,
		Quantity:     rQuantity,
		PricePerUnit: rShopPrice.PricePerUnit,
		Description:  rShopPrice.Description,
		Naming: maybeTypeNaming(
			syncNamingSession,
			rShopPrice.TypeId,
		),
	}
}
