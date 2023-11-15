package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/purchasequeue"
)

type PBShopPurchaseQueueParams struct{}

type PBShopPurchaseQueueClient struct{}

func NewPBShopPurchaseQueueClient() PBShopPurchaseQueueClient {
	return PBShopPurchaseQueueClient{}
}

func (gspqc PBShopPurchaseQueueClient) Fetch(
	x cache.Context,
	params PBShopPurchaseQueueParams,
) (
	entries []*proto.PurchaseQueueEntry,
	err error,
) {
	rShopQueue, _, err := purchasequeue.GetPurchaseQueue(x)
	if err != nil {
		return entries, err
	}

	entries = make([]*proto.PurchaseQueueEntry, 0)

	for _, appraisalCodes := range rShopQueue {
		for _, appraisalCode := range appraisalCodes {
			entries = append(
				entries,
				&proto.PurchaseQueueEntry{Code: appraisalCode},
			)
		}
	}

	return entries, nil
}
