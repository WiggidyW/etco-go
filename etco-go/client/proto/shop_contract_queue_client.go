package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/structureinfo"
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

func NewPBShopContractQueueClient(
	rContractsClient contracts.WC_ContractsClient,
	structureInfoClient structureinfo.WC_StructureInfoClient,
) PBShopContractQueueClient {
	return PBShopContractQueueClient{
		innerClient: NewPBContractQueueClient(
			rContractsClient,
			structureInfoClient,
		),
	}
}

func (scqc PBShopContractQueueClient) Fetch(
	ctx context.Context,
	params PBShopContractQueueParams,
) (
	entries []*proto.ContractQueueEntry,
	err error,
) {
	return scqc.innerClient.Fetch(
		ctx,
		PBContractQueueParams{
			LocationInfoSession: params.LocationInfoSession,
			StoreKind:           kind.Shop,
		},
	)
}
