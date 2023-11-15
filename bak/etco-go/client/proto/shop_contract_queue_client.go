package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/kind"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBShopContractQueueParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
}

type PBShopContractQueueClient struct {
	innerClient PBContractQueueClient
}

func NewPBShopContractQueueClient() PBShopContractQueueClient {
	return PBShopContractQueueClient{
		innerClient: NewPBContractQueueClient(),
	}
}

func (scqc PBShopContractQueueClient) Fetch(
	x cache.Context,
	params PBShopContractQueueParams,
) (
	entries []*proto.ContractQueueEntry,
	err error,
) {
	return scqc.innerClient.Fetch(
		x,
		PBContractQueueParams{
			LocationInfoSession: params.LocationInfoSession,
			StoreKind:           kind.Shop,
		},
	)
}
