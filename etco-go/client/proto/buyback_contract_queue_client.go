package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/structureinfo"
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

func NewPBBuybackContractQueueClient(
	rContractsClient contracts.WC_ContractsClient,
	structureInfoClient structureinfo.WC_StructureInfoClient,
) PBBuybackContractQueueClient {
	return PBBuybackContractQueueClient{
		innerClient: NewPBContractQueueClient(
			rContractsClient,
			structureInfoClient,
		),
	}
}

func (scqc PBBuybackContractQueueClient) Fetch(
	ctx context.Context,
	params PBBuybackContractQueueParams,
) (
	entries []*proto.ContractQueueEntry,
	err error,
) {
	return scqc.innerClient.Fetch(
		ctx,
		PBContractQueueParams{
			LocationInfoSession: params.LocationInfoSession,
			StoreKind:           kind.Buyback,
		},
	)
}
