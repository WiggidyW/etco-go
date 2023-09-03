package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func NewPBShopAppraisal[IM staticdb.IndexMap](
	rAppraisal rdb.ShopAppraisal,
	namingSession *staticdb.TypeNamingSession[IM],
) *proto.ShopAppraisal {
	pbAppraisal := &proto.ShopAppraisal{
		Items: make(
			[]*proto.ShopItem,
			0,
			len(rAppraisal.Items),
		),
		Code:       "",
		Price:      rAppraisal.Price,
		Time:       rAppraisal.Time.Unix(),
		Version:    rAppraisal.Version,
		LocationId: rAppraisal.LocationId,
	}
	for _, rItem := range rAppraisal.Items {
		pbAppraisal.Items = append(
			pbAppraisal.Items,
			NewPBShopItem(rItem, namingSession),
		)
	}
	return pbAppraisal
}

func NewPBShopItem[IM staticdb.IndexMap](
	rShopItem rdb.ShopItem,
	namingSession *staticdb.TypeNamingSession[IM],
) *proto.ShopItem {
	return &proto.ShopItem{
		TypeId:       rShopItem.TypeId,
		Quantity:     rShopItem.Quantity,
		PricePerUnit: rShopItem.PricePerUnit,
		Description:  rShopItem.Description,
		TypeNamingIndexes: MaybeGetTypeNamingIndexes(
			namingSession,
			rShopItem.TypeId,
		),
	}
}
