package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/kind"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBBuybackContractQueueParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
}

type PBBuybackContractQueueClient struct {
	innerClient PBContractQueueClient
}

func NewPBBuybackContractQueueClient() PBBuybackContractQueueClient {
	return PBBuybackContractQueueClient{
		innerClient: NewPBContractQueueClient(),
	}
}

func (scqc PBBuybackContractQueueClient) Fetch(
	x cache.Context,
	params PBBuybackContractQueueParams,
) (
	entries []*proto.ContractQueueEntry,
	err error,
) {
	return scqc.innerClient.Fetch(
		x,
		PBContractQueueParams{
			LocationInfoSession: params.LocationInfoSession,
			StoreKind:           kind.Buyback,
		},
	)
}
