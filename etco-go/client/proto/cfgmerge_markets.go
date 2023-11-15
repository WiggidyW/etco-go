package proto

import (
	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
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

type CfgMergeMarketsClient struct{}

func NewCfgMergeMarketsClient() CfgMergeMarketsClient {
	return CfgMergeMarketsClient{}
}

func (mmc CfgMergeMarketsClient) Fetch(
	x cache.Context,
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

	x, cancel := x.WithCancel()
	defer cancel()

	chnSendMarkets, chnRecvMarkets := chanresult.
		NewChanResult[map[b.MarketName]b.WebMarket](x.Ctx(), 1, 0).Split()
	go mmc.transceiveFetchMarkets(x, chnSendMarkets)

	// check if active map names are needed
	var activeMarketsHashSet util.MapHashSet[string, struct{}]
	if activeMapNamesNeeded(params.Updates) {
		activeMarketsHashSet, err = mmc.fetchActiveMarketsHashSet(x)
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
	if err = mmc.fetchWriteUpdated(x, markets); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mmc CfgMergeMarketsClient) fetchWriteUpdated(
	x cache.Context,
	updated map[b.MarketName]b.WebMarket,
) error {
	return bucket.SetWebMarkets(x, updated)
}

func (mmc CfgMergeMarketsClient) transceiveFetchMarkets(
	x cache.Context,
	chnSend chanresult.ChanSendResult[map[b.MarketName]b.WebMarket],
) error {
	markets, err := mmc.fetchMarkets(x)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(markets)
	}
}

func (mmc CfgMergeMarketsClient) fetchMarkets(
	x cache.Context,
) (
	markets map[b.MarketName]b.WebMarket,
	err error,
) {
	markets, _, err = bucket.GetWebMarkets(x)
	return markets, err
}

func (mmc CfgMergeMarketsClient) fetchActiveMarketsHashSet(
	x cache.Context,
) (
	activeMarketsHashSet util.MapHashSet[string, struct{}],
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	chnSendActiveMarkets, chnRecvActiveMarkets := chanresult.
		NewChanResult[map[string]struct{}](x.Ctx(), 1, 0).Split()
	go mmc.transceiveFetchSBuilderActiveMarkets(x, chnSendActiveMarkets)

	bigActiveMarkets, err := mmc.fetchBBuilderActiveMarkets(x)
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

func (mmc CfgMergeMarketsClient) fetchBBuilderActiveMarkets(
	x cache.Context,
) (
	activeMarkets map[string]struct{},
	err error,
) {
	bBuilder, _, err := bucket.GetWebBuybackSystemTypeMapsBuilder(x)
	if err != nil {
		return nil, err
	}
	return extractBuybackBuilderActiveMarkets(bBuilder), nil
}

func (mmc CfgMergeMarketsClient) transceiveFetchSBuilderActiveMarkets(
	x cache.Context,
	chnSend chanresult.ChanSendResult[map[string]struct{}],
) error {
	activeMarkets, err := mmc.fetchSBuilderActiveMarkets(x)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(activeMarkets)
	}
}

func (mmc CfgMergeMarketsClient) fetchSBuilderActiveMarkets(
	x cache.Context,
) (
	activeMarkets map[string]struct{},
	err error,
) {
	sBuilder, _, err := bucket.GetWebShopLocationTypeMapsBuilder(x)
	if err != nil {
		return nil, err
	}
	return extractShopBuilderActiveMarkets(sBuilder), nil
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
		if pbMarket == nil || pbMarket.LocationId == 0 {
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
