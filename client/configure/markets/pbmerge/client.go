package pbmerge

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	a "github.com/WiggidyW/weve-esi/client/authingfwding/authing"
	cfg "github.com/WiggidyW/weve-esi/client/configure"
	bbuilderget "github.com/WiggidyW/weve-esi/client/configure/btypemapsbuilder/get"
	fkbucketwriter "github.com/WiggidyW/weve-esi/client/configure/internal/fixedkeybucket/writer"
	marketsget "github.com/WiggidyW/weve-esi/client/configure/markets/get"
	sbuilderget "github.com/WiggidyW/weve-esi/client/configure/stypemapsbuilder/get"
	"github.com/WiggidyW/weve-esi/util"
)

var EmptyHashSet = util.MapHashSet[string, struct{}](map[string]struct{}{})

type A_PbMergeMarketsClient = a.AuthingClient[
	authingfwding.WithAuthableParams[PbMergeMarketsParams],
	PbMergeMarketsParams,
	cfg.PBMergeResponse,
	PbMergeMarketsClient,
]

type PbMergeMarketsClient struct {
	GetSBuilderClient sbuilderget.GetShopLocationTypeMapsBuilderClient
	GetBBuilderClient bbuilderget.GetBuybackSystemTypeMapsBuilderClient
	GetMarketsClient  marketsget.GetMarketsClient
	WriteClient       fkbucketwriter.SAC_FixedKeyBucketWriterClient[cfg.Markets]
}

func (mbbc PbMergeMarketsClient) Fetch(
	ctx context.Context,
	params PbMergeMarketsParams,
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

	// fetch the original markets and active map names
	var markets cfg.Markets
	var activeMapNames util.MapHashSet[string, struct{}]
	var err error
	if ActiveMapNamesNeeded(params.Updates.Inner) {
		markets, activeMapNames, err = mbbc.fetchMergeDataActive(
			ctx,
			params,
		)
	} else {
		markets, activeMapNames, err = mbbc.fetchMergeDataEmpty(
			ctx,
			params,
		)
	}
	if err != nil {
		return nil, err
	}

	// mutate the original markets with the updates
	if err = MergeMarkets(
		markets,
		params.Updates.Inner,
		activeMapNames,
	); err != nil {
		return &cfg.PBMergeResponse{
			// Modified:   false,
			MergeError: err,
		}, nil
	}

	// write the mutated markets
	if err = mbbc.writeUpdated(ctx, markets); err != nil {
		return nil, err
	}

	return &cfg.PBMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mmc PbMergeMarketsClient) fetchMergeDataEmpty(
	ctx context.Context,
	params PbMergeMarketsParams,
) (cfg.Markets, util.MapHashSet[string, struct{}], error) {
	// fetch the original markets
	if markets, err := mmc.fetchMarkets(ctx); err != nil {
		return nil, nil, err
	} else {
		return markets, EmptyHashSet, nil
	}
}

func (mmc PbMergeMarketsClient) fetchMergeDataActive(
	ctx context.Context,
	params PbMergeMarketsParams,
) (cfg.Markets, util.MapHashSet[string, struct{}], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the active markets for both buyback and shop
	chnActiveMarkets := util.NewChanResult[map[string]struct{}](ctx)
	chnSendActiveMarkets, chnRecvActiveMarkets := chnActiveMarkets.Split()
	go mmc.fetchBBuilderActiveMarkets(ctx, chnSendActiveMarkets)
	go mmc.fetchSBuilderActiveMarkets(ctx, chnSendActiveMarkets)

	// fetch the original markets
	markets, err := mmc.fetchMarkets(ctx)
	if err != nil {
		return nil, nil, err
	}

	// wait for the active markets
	var bigActiveMarkets, smallActiveMarkets map[string]struct{}
	if bigActiveMarkets, err = chnRecvActiveMarkets.Recv(); err != nil {
		return nil, nil, err
	}
	if smallActiveMarkets, err = chnRecvActiveMarkets.Recv(); err != nil {
		return nil, nil, err
	}
	if len(smallActiveMarkets) > len(bigActiveMarkets) {
		bigActiveMarkets, smallActiveMarkets = smallActiveMarkets,
			bigActiveMarkets
	}

	// insert the small active markets into the big active markets
	for market := range smallActiveMarkets {
		bigActiveMarkets[market] = struct{}{}
	}

	return markets, util.MapHashSet[string, struct{}](bigActiveMarkets), nil
}

func (mmc PbMergeMarketsClient) fetchMarkets(
	ctx context.Context,
) (cfg.Markets, error) {
	if marketsRep, err := mmc.GetMarketsClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return nil, err
	} else {
		return marketsRep.Data(), nil
	}
}

func (mmc PbMergeMarketsClient) fetchBBuilderActiveMarkets(
	ctx context.Context,
	chnSendActiveMarkets util.ChanSendResult[map[string]struct{}],
) error {
	if bBuilderRep, err := mmc.GetBBuilderClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return chnSendActiveMarkets.SendErr(err)
	} else {
		bBuilder := bBuilderRep.Data()
		activeMarkets := extractBuybackBuilderActiveMarkets(bBuilder)
		return chnSendActiveMarkets.SendOk(activeMarkets)
	}
}

func (mmc PbMergeMarketsClient) fetchSBuilderActiveMarkets(
	ctx context.Context,
	chnSendActiveMarkets util.ChanSendResult[map[string]struct{}],
) error {
	if sBuilderRep, err := mmc.GetSBuilderClient.Fetch(
		ctx,
		struct{}{},
	); err != nil {
		return chnSendActiveMarkets.SendErr(err)
	} else {
		sBuilder := sBuilderRep.Data()
		activeMarkets := extractShopBuilderActiveMarkets(sBuilder)
		return chnSendActiveMarkets.SendOk(activeMarkets)
	}
}

func (mmc PbMergeMarketsClient) writeUpdated(
	ctx context.Context,
	updated cfg.Markets,
) error {
	if _, err := mmc.WriteClient.Fetch(ctx, fkbucketwriter.
		FixedKeyBucketWriterParams[cfg.Markets]{
		Val: updated,
	}); err != nil {
		return err
	} else {
		return nil
	}
}
