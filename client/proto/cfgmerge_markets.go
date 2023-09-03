package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
	"github.com/WiggidyW/etco-go/util"
)

var emptyActiveMarketsHashSet = util.MapHashSet[string, struct{}](
	map[string]struct{}{},
)

type CfgMergeMarketsParams struct {
	Updates map[string]*proto.CfgMarket
}

type CfgMergeMarketsClient struct {
	webMarketsReaderClient          bucket.WebMarketsReaderClient
	webMarketsWriterClient          bucket.WebMarketsWriterClient
	webBTypeMapsBuilderReaderClient bucket.WebBuybackSystemTypeMapsBuilderReaderClient
	webSTypeMapsBuilderReaderClient bucket.WebShopLocationTypeMapsBuilderReaderClient
}

func (mmc CfgMergeMarketsClient) Fetch(
	ctx context.Context,
	params CfgMergeMarketsParams,
) (
	rep *CfgMergeResponse,
	err error,
) {
	// if there are no updates, return now
	if params.Updates == nil || len(params.Updates) == 0 {
		return &CfgMergeResponse{
			// Modified: false,
			// MergeError: nil,
		}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendMarkets, chnRecvMarkets := chanresult.
		NewChanResult[map[b.MarketName]b.WebMarket](ctx, 1, 0).Split()
	go mmc.transceiveFetchMarkets(ctx, chnSendMarkets)

	// check if active map names are needed
	var activeMarketsHashSet util.MapHashSet[string, struct{}]
	if activeMapNamesNeeded(params.Updates) {
		activeMarketsHashSet, err = mmc.fetchActiveMarketsHashSet(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		activeMarketsHashSet = emptyActiveMarketsHashSet
	}

	// receive the markets
	markets, err := chnRecvMarkets.Recv()
	if err != nil {
		return nil, err
	}

	// mutate the original markets with the updates
	if err = mergeMarkets(
		markets,
		params.Updates,
		activeMarketsHashSet,
	); err != nil {
		return &CfgMergeResponse{
			// Modified:   false,
			MergeError: err,
		}, nil
	}

	// write the mutated markets
	if err = mmc.fetchWriteUpdated(ctx, markets); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mmc CfgMergeMarketsClient) fetchWriteUpdated(
	ctx context.Context,
	updated map[b.MarketName]b.WebMarket,
) error {
	_, err := mmc.webMarketsWriterClient.Fetch(
		ctx,
		bucket.WebMarketsWriterParams{
			WebMarkets: updated,
		},
	)
	return err
}

func (mmc CfgMergeMarketsClient) transceiveFetchMarkets(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[map[b.MarketName]b.WebMarket],
) error {
	markets, err := mmc.fetchMarkets(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(markets)
	}
}

func (mmc CfgMergeMarketsClient) fetchMarkets(
	ctx context.Context,
) (
	markets map[b.MarketName]b.WebMarket,
	err error,
) {
	marketsRep, err := mmc.webMarketsReaderClient.Fetch(
		ctx,
		bucket.WebMarketsReaderParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return marketsRep.Data(), nil
	}
}

// func (mmc CfgMergeMarketsClient) transceiveFetchActiveMarketsHashSet(
// 	ctx context.Context,
// 	chnSend chanresult.ChanSendResult[util.MapHashSet[string, struct{}]],
// ) error {
// 	activeMarketsHashSet, err := mmc.fetchActiveMarketsHashSet(ctx)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(activeMarketsHashSet)
// 	}
// }

func (mmc CfgMergeMarketsClient) fetchActiveMarketsHashSet(
	ctx context.Context,
) (
	activeMarketsHashSet util.MapHashSet[string, struct{}],
	err error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendActiveMarkets, chnRecvActiveMarkets := chanresult.
		NewChanResult[map[string]struct{}](ctx, 1, 0).Split()
	go mmc.transceiveFetchSBuilderActiveMarkets(ctx, chnSendActiveMarkets)

	bigActiveMarkets, err := mmc.fetchBBuilderActiveMarkets(ctx)
	if err != nil {
		return nil, err
	}

	smallActiveMarkets, err := chnRecvActiveMarkets.Recv()
	if err != nil {
		return nil, err
	}

	if len(smallActiveMarkets) > len(bigActiveMarkets) {
		bigActiveMarkets, smallActiveMarkets = smallActiveMarkets,
			bigActiveMarkets
	}

	// insert the small market names into the big market names
	for activeMarket := range smallActiveMarkets {
		bigActiveMarkets[activeMarket] = struct{}{}
	}

	return util.MapHashSet[string, struct{}](bigActiveMarkets), nil
}

// func (mmc CfgMergeMarketsClient) transceiveFetchBBuilderActiveMarkets(
// 	ctx context.Context,
// 	chnSend chanresult.ChanSendResult[map[string]struct{}],
// ) error {
// 	activeMarkets, err := mmc.fetchBBuilderActiveMarkets(ctx)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(activeMarkets)
// 	}
// }

func (mmc CfgMergeMarketsClient) fetchBBuilderActiveMarkets(
	ctx context.Context,
) (
	activeMarkets map[string]struct{},
	err error,
) {
	bBuilderRep, err := mmc.webBTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebBuybackSystemTypeMapsBuilderReaderParams{},
	)
	if err != nil {
		return nil, err
	}
	return extractBuybackBuilderActiveMarkets(bBuilderRep.Data()), nil
}

func (mmc CfgMergeMarketsClient) transceiveFetchSBuilderActiveMarkets(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[map[string]struct{}],
) error {
	activeMarkets, err := mmc.fetchSBuilderActiveMarkets(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(activeMarkets)
	}
}

func (mmc CfgMergeMarketsClient) fetchSBuilderActiveMarkets(
	ctx context.Context,
) (
	activeMarkets map[string]struct{},
	err error,
) {
	sBuilderRep, err := mmc.webSTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebShopLocationTypeMapsBuilderReaderParams{},
	)
	if err != nil {
		return nil, err
	}
	return extractShopBuilderActiveMarkets(sBuilderRep.Data()), nil
}

func newPBtoWebMarketError(
	name string,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{Err: configerror.ErrMarketInvalid{
		Market:    name,
		ErrString: errStr,
	}}
}

func pBtoWebMarket(
	marketName string,
	pbMarket *proto.CfgMarket,
) (
	webMarket *b.WebMarket,
	err error,
) {
	if pbMarket.IsStructure && pbMarket.RefreshToken == "" {
		return nil, newPBtoWebMarketError(
			marketName,
			"structure market must have refresh token",
		)
	} else if !pbMarket.IsStructure {
		if pbMarket.RefreshToken != "" {
			return nil, newPBtoWebMarketError(
				marketName,
				"non-structure market must not have refresh token",
			)
		} else if stationInfo := staticdb.GetStationInfo(
			int32(pbMarket.LocationId),
		); stationInfo == nil {
			return nil, newPBtoWebMarketError(
				marketName,
				"station does not exist",
			)
		}
	}

	webMarket = &b.WebMarket{
		LocationId:  pbMarket.LocationId,
		IsStructure: pbMarket.IsStructure,
	}
	if pbMarket.RefreshToken != "" {
		webMarket.RefreshToken = &pbMarket.RefreshToken
	}

	return webMarket, nil
}

func mergeMarkets[HS util.HashSet[string]](
	original map[b.MarketName]b.WebMarket,
	updates map[string]*proto.CfgMarket,
	activeMapNames HS,
) error {
	for marketName, pbMarket := range updates {
		if pbMarket == nil {
			if activeMapNames.Has(marketName) {
				return newPBtoWebMarketError(
					marketName,
					"cannot delete: market currently in use",
				)
			} else {
				delete(original, marketName)
			}
		} else {
			webMarket, err := pBtoWebMarket(marketName, pbMarket)
			if err != nil {
				return err
			}
			original[marketName] = *webMarket
		}
	}
	return nil
}

// since getting ActiveMapNames is expensive operation, check if it's needed
func activeMapNamesNeeded(updates map[string]*proto.CfgMarket) bool {
	for _, pbMarket := range updates {
		if pbMarket == nil {
			return true
		}
	}
	return false
}
