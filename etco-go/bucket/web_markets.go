package bucket

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/staticdb"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_MARKETS_BUF_CAP          int           = 0
	WEB_MARKETS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_MARKETS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_MARKETS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebMarkets = cache.RegisterType[map[b.MarketName]b.WebMarket]("webmarkets", WEB_MARKETS_BUF_CAP)
}

func GetWebMarkets(
	x cache.Context,
) (
	rep map[b.MarketName]b.WebMarket,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebMarkets,
		keys.CacheKeyWebMarkets,
		keys.TypeStrWebMarkets,
		WEB_MARKETS_EXPIRES_IN,
		build.CAPACITY_WEB_MARKETS,
	)
}

func ProtoGetWebMarketNames(
	x cache.Context,
) (
	rep []string,
	expires time.Time,
	err error,
) {
	var webMarkets map[b.MarketName]b.WebMarket
	webMarkets, expires, err = GetWebMarkets(x)
	if err == nil {
		rep = keysToSlice(webMarkets)
	}
	return rep, expires, err
}

func ProtoGetWebMarkets(
	x cache.Context,
) (
	rep map[string]*proto.CfgMarket,
	expires time.Time,
	err error,
) {
	var webMarkets map[b.MarketName]b.WebMarket
	webMarkets, expires, err = GetWebMarkets(x)
	if err == nil {
		rep = WebMarketsToProto(webMarkets)
	}
	return rep, expires, err
}

func SetWebMarkets(
	x cache.Context,
	rep map[b.MarketName]b.WebMarket,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebMarkets,
		keys.CacheKeyWebMarkets,
		keys.TypeStrWebMarkets,
		WEB_MARKETS_EXPIRES_IN,
		rep,
		nil,
	)
}

func ProtoMergeSetWebMarkets(
	x cache.Context,
	updates map[string]*proto.CfgMarket,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}

	// fetch the original markets in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chn := expirable.NewChanResult[map[b.MarketName]b.WebMarket](x.Ctx(), 1, 0)
	go expirable.P1Transceive(chn, x, GetWebMarkets)

	// check if active markets are needed and get them if so
	var activeMarkets map[string]struct{}
	if activeMapNamesNeeded(updates) {
		activeMarkets, _, err = GetWebActiveMarkets(x)
		if err != nil {
			return err
		}
	}

	// receive the original markets
	var markets map[b.MarketName]b.WebMarket
	markets, _, err = chn.RecvExp()
	if err != nil {
		return err
	}

	// merge updates
	err = ProtoMergeMarkets(markets, updates, activeMarkets)
	if err != nil {
		return protoerr.New(protoerr.INVALID_MERGE, err)
	}

	return SetWebMarkets(x, markets)
}

// // To Proto

func WebMarketsToProto(
	webMarkets map[b.MarketName]b.WebMarket,
) (
	pbMarkets map[string]*proto.CfgMarket,
) {
	return newPBCfgMarkets(webMarkets)
}

func newPBCfgMarkets(
	webMarkets map[b.MarketName]b.WebMarket,
) (
	pbMarkets map[string]*proto.CfgMarket,
) {
	pbMarkets = make(
		map[string]*proto.CfgMarket,
		len(webMarkets),
	)
	for marketName, webMarket := range webMarkets {
		pbMarkets[marketName] = newPBCfgMarket(webMarket)
	}
	return pbMarkets
}

func newPBCfgMarket(
	webMarket b.WebMarket,
) (
	pbMarket *proto.CfgMarket,
) {
	if webMarket.RefreshToken != nil {
		return &proto.CfgMarket{
			RefreshToken: *webMarket.RefreshToken,
			LocationId:   webMarket.LocationId,
			IsStructure:  webMarket.IsStructure,
		}
	} else {
		return &proto.CfgMarket{
			// RefreshToken: "",
			LocationId:  webMarket.LocationId,
			IsStructure: webMarket.IsStructure,
		}
	}
}

// // Merge

func ProtoMergeMarkets(
	original map[b.MarketName]b.WebMarket,
	updates map[string]*proto.CfgMarket,
	activeMapNames map[string]struct{},
) error {
	for marketName, pbMarket := range updates {
		if pbMarket == nil || pbMarket.LocationId == 0 {
			if _, ok := activeMapNames[marketName]; ok {
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
