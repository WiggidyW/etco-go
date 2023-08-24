package pbmerge

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
	a "github.com/WiggidyW/eve-trading-co-go/client/authingfwding/authing"
	cfg "github.com/WiggidyW/eve-trading-co-go/client/configure"
	fkbucketwriter "github.com/WiggidyW/eve-trading-co-go/client/configure/internal/fixedkeybucket/writer"
	marketsget "github.com/WiggidyW/eve-trading-co-go/client/configure/markets/get"
	builderget "github.com/WiggidyW/eve-trading-co-go/client/configure/stypemapsbuilder/get"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type A_PbMergeShopLocationTypeMapsBuilderClient = a.AuthingClient[
	authingfwding.WithAuthableParams[PbMergeShopLocationTypeMapsBuilderParams],
	PbMergeShopLocationTypeMapsBuilderParams,
	cfg.PBMergeResponse,
	PbMergeShopLocationTypeMapsBuilderClient,
]

// TODO: hideous name
type PbMergeShopLocationTypeMapsBuilderClient struct {
	GetBuilderClient builderget.GetShopLocationTypeMapsBuilderClient
	GetMarketsClient marketsget.GetMarketsClient
	WriteClient      fkbucketwriter.SAC_FixedKeyBucketWriterClient[cfg.ShopLocationTypeMapsBuilder]
}

func (msbc PbMergeShopLocationTypeMapsBuilderClient) Fetch(
	ctx context.Context,
	params PbMergeShopLocationTypeMapsBuilderParams,
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the original builder
	chnBuilder := util.NewChanResult[cfg.ShopLocationTypeMapsBuilder](ctx)
	chnBuilderSend, chnBuilderRecv := chnBuilder.Split()
	go msbc.fetchBuilder(ctx, chnBuilderSend)

	// fetch markets (used for update validation - ensures markets exist)
	marketHashSet, err := msbc.fetchMarketsHashSet(ctx)
	if err != nil {
		return nil, err
	}

	// wait for the original builder
	builder, err := chnBuilderRecv.Recv()
	if err != nil {
		return nil, err
	}

	// mutate the original builder with the updates
	if err := MergeShopLocationTypeMapsBuilder(
		builder,
		params.Updates.Inner,
		marketHashSet,
	); err != nil {
		// return the error
		return &cfg.PBMergeResponse{
			// Modified:   false,
			MergeError: err,
		}, nil
	}

	if err := msbc.writeUpdated(ctx, builder); err != nil {
		return nil, err
	}

	return &cfg.PBMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (msbc PbMergeShopLocationTypeMapsBuilderClient) writeUpdated(
	ctx context.Context,
	updated cfg.ShopLocationTypeMapsBuilder,
) error {
	if _, err := msbc.WriteClient.Fetch(ctx, fkbucketwriter.
		FixedKeyBucketWriterParams[cfg.ShopLocationTypeMapsBuilder]{
		Val: updated,
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (msbc PbMergeShopLocationTypeMapsBuilderClient) fetchMarketsHashSet(
	ctx context.Context,
) (util.MapHashSet[string, cfg.Market], error) {
	if marketsRep, err := msbc.GetMarketsClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return nil, err
	} else {
		markets := marketsRep.Data()
		return util.MapHashSet[string, cfg.Market](markets), nil
	}
}

func (msbc PbMergeShopLocationTypeMapsBuilderClient) fetchBuilder(
	ctx context.Context,
	chnBuilderSend util.ChanSendResult[cfg.ShopLocationTypeMapsBuilder],
) error {
	if builderRep, err := msbc.GetBuilderClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return chnBuilderSend.SendErr(err)
	} else {
		return chnBuilderSend.SendOk(builderRep.Data())
	}
}
