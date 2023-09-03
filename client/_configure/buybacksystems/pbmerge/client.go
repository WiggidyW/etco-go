package pbmerge

import (
	"context"

	"github.com/WiggidyW/chanresult"
	"github.com/WiggidyW/etco-go/client/authingfwding"
	a "github.com/WiggidyW/etco-go/client/authingfwding/authing"
	cfg "github.com/WiggidyW/etco-go/client/configure"
	bbuilderget "github.com/WiggidyW/etco-go/client/configure/btypemapsbuilder/get"
	systemsget "github.com/WiggidyW/etco-go/client/configure/buybacksystems/get"
	fkbucketwriter "github.com/WiggidyW/etco-go/client/configure/internal/fixedkeybucket/writer"
	sbuilderget "github.com/WiggidyW/etco-go/client/configure/stypemapsbuilder/get"
	"github.com/WiggidyW/etco-go/util"
)

type A_PbMergeBuybackSystemsBuilderClient = a.AuthingClient[
	authingfwding.WithAuthableParams[PbMergeBuybackSystemsParams],
	PbMergeBuybackSystemsParams,
	cfg.PBMergeResponse,
	PbMergeBuybackSystemsClient,
]

type PbMergeBuybackSystemsClient struct {
	GetSBuilderClient sbuilderget.GetShopLocationTypeMapsBuilderClient
	GetBBuilderClient bbuilderget.GetBuybackSystemTypeMapsBuilderClient
	GetSystemsClient  systemsget.GetBuybackSystemsClient
	WriteClient       fkbucketwriter.SAC_FixedKeyBucketWriterClient[cfg.BuybackSystems]
}

func (mbsc PbMergeBuybackSystemsClient) Fetch(
	ctx context.Context,
	params PbMergeBuybackSystemsParams,
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
	if err = MergeBuybackSystems(
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

func (mbsc PbMergeBuybackSystemsClient) fetchMergeData(
	ctx context.Context,
	params PbMergeBuybackSystemsParams,
) (cfg.BuybackSystems, util.MapHashSet[string, struct{}], error) {
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

func (mbsc PbMergeBuybackSystemsClient) fetchSystems(
	ctx context.Context,
) (cfg.BuybackSystems, error) {
	if systemsRep, err := mbsc.GetSystemsClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return nil, err
	} else {
		return systemsRep.Data(), nil
	}
}

func (mbsc PbMergeBuybackSystemsClient) fetchBBuilderBundleKeys(
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

func (mbsc PbMergeBuybackSystemsClient) fetchSBuilderBundleKeys(
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

func (mbsc PbMergeBuybackSystemsClient) writeUpdated(
	ctx context.Context,
	updated cfg.BuybackSystems,
) error {
	if _, err := mbsc.WriteClient.Fetch(ctx, fkbucketwriter.
		FixedKeyBucketWriterParams[cfg.BuybackSystems]{
		Val: updated,
	}); err != nil {
		return err
	} else {
		return nil
	}
}
