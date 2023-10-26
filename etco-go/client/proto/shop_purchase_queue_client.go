package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/shopqueue"
	"github.com/WiggidyW/etco-go/proto"
)

type PBShopPurchaseQueueParams struct{}

type PBShopPurchaseQueueClient struct {
	rShopQueueClient shopqueue.ShopQueueClient
}

func NewPBShopPurchaseQueueClient(
	rShopQueueClient shopqueue.ShopQueueClient,
) PBShopPurchaseQueueClient {
	return PBShopPurchaseQueueClient{rShopQueueClient}
}

func (gspqc PBShopPurchaseQueueClient) Fetch(
	ctx context.Context,
	params PBShopPurchaseQueueParams,
) (
	entries []*proto.PurchaseQueueEntry,
	err error,
) {
	rShopQueueRep, err := gspqc.rShopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	)
	if err != nil {
		return entries, err
	}
	rShopQueue := rShopQueueRep.ParsedShopQueue

	entries = make(
		[]*proto.PurchaseQueueEntry,
		0,
		len(rShopQueue),
	)

	for _, appraisalCode := range rShopQueue {
		entries = append(
			entries,
			&proto.PurchaseQueueEntry{Code: appraisalCode},
		)
	}

	return entries, nil
}
