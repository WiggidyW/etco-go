package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/proto/protoerr"

	b "github.com/WiggidyW/etco-go-bucket"
)

func GetWebActiveMarkets(
	x cache.Context,
) (
	rep map[string]struct{},
	expires time.Time,
	err error,
) {
	// fetch shop active markets in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chn := expirable.NewChanResult[map[string]struct{}](x.Ctx(), 1, 0)
	go expirable.P1Transceive(chn, x, GetWebShopLocationActiveMarkets)

	// fetch buyback active markets
	var biggerActiveMarkets map[string]struct{}
	biggerActiveMarkets, expires, err = GetWebBuybackSystemActiveMarkets(x)
	if err != nil {
		return nil, expires, err
	}

	// recv shop active markets
	var smallerActiveMarkets map[string]struct{}
	smallerActiveMarkets, expires, err = chn.RecvExpMin(expires)
	if err != nil {
		return nil, expires, err
	}

	// merge active markets
	if len(biggerActiveMarkets) < len(smallerActiveMarkets) {
		biggerActiveMarkets, smallerActiveMarkets = smallerActiveMarkets,
			biggerActiveMarkets
	}
	for market := range smallerActiveMarkets {
		biggerActiveMarkets[market] = struct{}{}
	}

	return biggerActiveMarkets, expires, nil
}

func protoMergeSetTypeMapsBuilder[U any, B any](
	x cache.Context,
	updates U,
	getOriginal func(cache.Context) (B, time.Time, error),
	mergeUpdates func(B, U, map[b.MarketName]b.WebMarket) error,
	setUpdated func(cache.Context, B) error,
) (
	err error,
) {
	// fetch markets in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnMarkets :=
		expirable.NewChanResult[map[b.MarketName]b.WebMarket](x.Ctx(), 1, 0)
	go expirable.P1Transceive(chnMarkets, x, GetWebMarkets)

	// fetch the original builder
	var builder B
	builder, _, err = getOriginal(x)
	if err != nil {
		return err
	}

	// recv markets hashset
	var markets map[b.MarketName]b.WebMarket
	markets, _, err = chnMarkets.RecvExp()
	if err != nil {
		return err
	}

	// merge updates
	err = mergeUpdates(builder, updates, markets)
	if err != nil {
		return protoerr.New(protoerr.INVALID_MERGE, err)
	}

	// set updated
	return setUpdated(x, builder)
}
