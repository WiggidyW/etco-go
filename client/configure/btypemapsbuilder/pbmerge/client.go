package pbmerge

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
	a "github.com/WiggidyW/eve-trading-co-go/client/authingfwding/authing"
	cfg "github.com/WiggidyW/eve-trading-co-go/client/configure"
	builderget "github.com/WiggidyW/eve-trading-co-go/client/configure/btypemapsbuilder/get"
	fkbucketwriter "github.com/WiggidyW/eve-trading-co-go/client/configure/internal/fixedkeybucket/writer"
	marketsget "github.com/WiggidyW/eve-trading-co-go/client/configure/markets/get"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type A_PbMergeBuybackSystemTypeMapsBuilderClient = a.AuthingClient[
	authingfwding.WithAuthableParams[PbMergeBuybackSystemTypeMapsBuilderParams],
	PbMergeBuybackSystemTypeMapsBuilderParams,
	cfg.PBMergeResponse,
	PbMergeBuybackSystemTypeMapsBuilderClient,
]

// TODO: hideous name
type PbMergeBuybackSystemTypeMapsBuilderClient struct {
	GetBuilderClient builderget.GetBuybackSystemTypeMapsBuilderClient
	GetMarketsClient marketsget.GetMarketsClient
	WriteClient      fkbucketwriter.SAC_FixedKeyBucketWriterClient[cfg.BuybackSystemTypeMapsBuilder]
}

func (mbbc PbMergeBuybackSystemTypeMapsBuilderClient) Fetch(
	ctx context.Context,
	params PbMergeBuybackSystemTypeMapsBuilderParams,
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
	chnBuilder := util.NewChanResult[cfg.BuybackSystemTypeMapsBuilder](ctx)
	chnBuilderSend, chnBuilderRecv := chnBuilder.Split()
	go mbbc.fetchBuilder(ctx, chnBuilderSend)

	// fetch markets (used for update validation - ensures markets exist)
	marketHashSet, err := mbbc.fetchMarketsHashSet(ctx)
	if err != nil {
		return nil, err
	}

	// wait for the original builder
	builder, err := chnBuilderRecv.Recv()
	if err != nil {
		return nil, err
	}

	// mutate the original builder with the updates
	if err := MergeBuybackSystemTypeMapsBuilder(
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

	if err := mbbc.writeUpdated(ctx, builder); err != nil {
		return nil, err
	}

	return &cfg.PBMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mbbc PbMergeBuybackSystemTypeMapsBuilderClient) writeUpdated(
	ctx context.Context,
	updated cfg.BuybackSystemTypeMapsBuilder,
) error {
	if _, err := mbbc.WriteClient.Fetch(ctx, fkbucketwriter.
		FixedKeyBucketWriterParams[cfg.BuybackSystemTypeMapsBuilder]{
		Val: updated,
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (mbbc PbMergeBuybackSystemTypeMapsBuilderClient) fetchMarketsHashSet(
	ctx context.Context,
) (util.MapHashSet[string, cfg.Market], error) {
	if marketsRep, err := mbbc.GetMarketsClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return nil, err
	} else {
		markets := marketsRep.Data()
		return util.MapHashSet[string, cfg.Market](markets), nil
	}
}

func (mbbc PbMergeBuybackSystemTypeMapsBuilderClient) fetchBuilder(
	ctx context.Context,
	chnBuilderSend util.ChanSendResult[cfg.BuybackSystemTypeMapsBuilder],
) error {
	if builderRep, err := mbbc.GetBuilderClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return chnBuilderSend.SendErr(err)
	} else {
		return chnBuilderSend.SendOk(builderRep.Data())
	}
}
