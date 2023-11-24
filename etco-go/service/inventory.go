package service

import (
	"context"

	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/shopassets"
)

func (Service) ShopInventory(
	ctx context.Context,
	req *proto.ShopInventoryRequest,
) (
	rep *proto.ShopInventoryResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.ShopInventoryResponse{}

	var inventoryRep shopassets.ProtoShopInventoryRep
	rep.Authorized, inventoryRep, err = authorizedGetP2(
		x,
		req.RefreshToken,
		auth.ProtoBoolUserAuthorized,
		shopassets.ProtoGetShopInventory,
		r,
		req.LocationId,
	)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
		rep.Strs = r.Finish()
		return rep, nil
	}

	rep.LocationInfo = inventoryRep.LocationInfo
	rep.Items = inventoryRep.Assets
	rep.Strs = r.Finish()
	return rep, nil
}
