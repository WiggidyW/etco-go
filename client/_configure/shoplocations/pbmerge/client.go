package pbmerge

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/authingfwding"
	a "github.com/WiggidyW/etco-go/client/authingfwding/authing"
	cfg "github.com/WiggidyW/etco-go/client/configure"
	bbuilderget "github.com/WiggidyW/etco-go/client/configure/btypemapsbuilder/get"
	fkbucketwriter "github.com/WiggidyW/etco-go/client/configure/internal/fixedkeybucket/writer"
	locationsget "github.com/WiggidyW/etco-go/client/configure/shoplocations/get"
	sbuilderget "github.com/WiggidyW/etco-go/client/configure/stypemapsbuilder/get"
	"github.com/WiggidyW/etco-go/util"
)

type A_PbMergeShopLocationsClient = a.AuthingClient[
	authingfwding.WithAuthableParams[PbMergeShopLocationsParams],
	PbMergeShopLocationsParams,
	cfg.PBMergeResponse,
	PbMergeShopLocationsClient,
]

type PbMergeShopLocationsClient struct {
	GetSBuilderClient sbuilderget.GetShopLocationTypeMapsBuilderClient
	GetBBuilderClient bbuilderget.GetBuybackSystemTypeMapsBuilderClient
	GetSystemsClient  locationsget.GetShopLocationsClient
	WriteClient       fkbucketwriter.SAC_FixedKeyBucketWriterClient[cfg.ShopLocations]
}

func (mbsc PbMergeShopLocationsClient) Fetch(
	ctx context.Context,
	params PbMergeShopLocationsParams,
) (*cfg.PBMergeResponse, error) {
	// if there are no updates, return now
	if params.Updates == nil ||
		params.Updates.Inner == nil ||
		len(params.Updates.Inner) == 0 {
		return &cfg.PBMergeResponse{
			// Modified: false,
			// MergeError: nil,
		}, nil
	}

	// fetch the original systems and active bundle keys
	systems, activeBundleKeys, err := mbsc.fetchMergeData(ctx, params)
	if err != nil {
		return nil, err
	}

	// mutate the original systems with the updates
	if err = MergeShopLocations(
		systems,
		params.Updates.Inner,
		activeBundleKeys,
	); err != nil {
		return &cfg.PBMergeResponse{
			// Modified: false,
			MergeError: err,
		}, nil
	}

	// write the mutated systems
	if err = mbsc.writeUpdated(ctx, systems); err != nil {
		return nil, err
	}

	return &cfg.PBMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mbsc PbMergeShopLocationsClient) fetchMergeData(
	ctx context.Context,
	params PbMergeShopLocationsParams,
) (cfg.ShopLocations, util.MapHashSet[string, struct{}], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the active bundle keys for both buyback and shop
	chnSendBundleKeys, chnRecvBundleKeys := chanresult.
		NewChanResult[map[string]struct{}](ctx, 0, 0).Split()
	go mbsc.fetchBBuilderBundleKeys(ctx, chnSendBundleKeys)
	go mbsc.fetchSBuilderBundleKeys(ctx, chnSendBundleKeys)

	// fetch the original systems
	systems, err := mbsc.fetchSystems(ctx)
	if err != nil {
		return nil, nil, err
	}

	// wait for the bundle keys
	var bigBundleKeys, smallBundleKeys map[string]struct{}
	if bigBundleKeys, err = chnRecvBundleKeys.Recv(); err != nil {
		return nil, nil, err
	}
	if smallBundleKeys, err = chnRecvBundleKeys.Recv(); err != nil {
		return nil, nil, err
	}
	if len(smallBundleKeys) > len(bigBundleKeys) {
		bigBundleKeys, smallBundleKeys = smallBundleKeys, bigBundleKeys
	}

	// insert the small bundle keys into the big bundle keys
	for bundleKey := range smallBundleKeys {
		bigBundleKeys[bundleKey] = struct{}{}
	}

	return systems, util.MapHashSet[string, struct{}](bigBundleKeys), nil
}

func (mbsc PbMergeShopLocationsClient) fetchSystems(
	ctx context.Context,
) (cfg.ShopLocations, error) {
	if systemsRep, err := mbsc.GetSystemsClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return nil, err
	} else {
		return systemsRep.Data(), nil
	}
}

func (mbsc PbMergeShopLocationsClient) fetchBBuilderBundleKeys(
	ctx context.Context,
	chnSendBundleKeys chanresult.ChanSendResult[map[string]struct{}],
) error {
	if bBuilderRep, err := mbsc.GetBBuilderClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return chnSendBundleKeys.SendErr(err)
	} else {
		bBuilder := bBuilderRep.Data()
		activeMarkets := cfg.ExtractBuilderBundleKeys(bBuilder)
		return chnSendBundleKeys.SendOk(activeMarkets)
	}
}

func (mbsc PbMergeShopLocationsClient) fetchSBuilderBundleKeys(
	ctx context.Context,
	chnSendBundleKeys chanresult.ChanSendResult[map[string]struct{}],
) error {
	if sBuilderRep, err := mbsc.GetSBuilderClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return chnSendBundleKeys.SendErr(err)
	} else {
		sBuilder := sBuilderRep.Data()
		activeMarkets := cfg.ExtractBuilderBundleKeys(sBuilder)
		return chnSendBundleKeys.SendOk(activeMarkets)
	}
}

func (mbsc PbMergeShopLocationsClient) writeUpdated(
	ctx context.Context,
	updated cfg.ShopLocations,
) error {
	if _, err := mbsc.WriteClient.Fetch(ctx, fkbucketwriter.
		FixedKeyBucketWriterParams[cfg.ShopLocations]{
		Val: updated,
	}); err != nil {
		return err
	} else {
		return nil
	}
}
