package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func NewPBBuybackAppraisal[IM staticdb.IndexMap](
	rAppraisal rdb.BuybackAppraisal,
	namingSession *staticdb.TypeNamingSession[IM],
) *proto.BuybackAppraisal {
	pbAppraisal := &proto.BuybackAppraisal{
		Items: make(
			[]*proto.BuybackParentItem,
			0,
			len(rAppraisal.Items),
		),
		Code:     rAppraisal.Code,
		Price:    rAppraisal.Price,
		Fee:      rAppraisal.Fee,
		Time:     rAppraisal.Time.Unix(),
		Version:  rAppraisal.Version,
		SystemId: rAppraisal.SystemId,
	}
	for _, rParentItem := range rAppraisal.Items {
		pbAppraisal.Items = append(
			pbAppraisal.Items,
			NewPBBuybackParentItem(
				rParentItem,
				namingSession,
			),
		)
	}
	return pbAppraisal
}

func NewPBBuybackParentItem[IM staticdb.IndexMap](
	rParentItem rdb.BuybackParentItem,
	namingSession *staticdb.TypeNamingSession[IM],
) *proto.BuybackParentItem {
	pbParentItem := &proto.BuybackParentItem{
		TypeId:       rParentItem.TypeId,
		Quantity:     rParentItem.Quantity,
		PricePerUnit: rParentItem.PricePerUnit,
		Fee:          rParentItem.Fee,
		Description:  rParentItem.Description,
		Children: make(
			[]*proto.BuybackChildItem,
			0,
			len(rParentItem.Children),
		),
		TypeNamingIndexes: MaybeGetTypeNamingIndexes(
			namingSession,
			rParentItem.TypeId,
		),
	}
	for _, rChildItem := range rParentItem.Children {
		pbParentItem.Children = append(
			pbParentItem.Children,
			NewPBBuybackChildItem(rChildItem, namingSession),
		)
	}
	return pbParentItem
}

func NewPBBuybackChildItem[IM staticdb.IndexMap](
	rChildItem rdb.BuybackChildItem,
	namingSession *staticdb.TypeNamingSession[IM],
) *proto.BuybackChildItem {
	return &proto.BuybackChildItem{
		TypeId:            rChildItem.TypeId,
		QuantityPerParent: rChildItem.QuantityPerParent,
		PricePerUnit:      rChildItem.PricePerUnit,
		Description:       rChildItem.Description,
		TypeNamingIndexes: MaybeGetTypeNamingIndexes(
			namingSession,
			rChildItem.TypeId,
		),
	}
}
