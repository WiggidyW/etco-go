package webtocore

import (
	"context"
	"time"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_BUCKET_DATA_TIMEOUT = 300 * time.Second
)

func downloadWebBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (b.WebBucketData, error) {
	ctx, cancel := context.WithTimeout(ctx, WEB_BUCKET_DATA_TIMEOUT)
	defer cancel()

	chnWebBSTypeMapsBuilder :=
		chanresult.NewChanResult[map[b.TypeId]b.WebBuybackSystemTypeBundle](
			ctx, 1, 0,
		)
	go transceiveFetchBucketData(
		ctx,
		chnWebBSTypeMapsBuilder.ToSend(),
		bucketClient.ReadWebBuybackSystemTypeMapsBuilder,
	)

	chnWebSLTypeMapsBuilder :=
		chanresult.NewChanResult[map[b.TypeId]b.WebShopLocationTypeBundle](
			ctx, 1, 0,
		)
	go transceiveFetchBucketData(
		ctx,
		chnWebSLTypeMapsBuilder.ToSend(),
		bucketClient.ReadWebShopLocationTypeMapsBuilder,
	)

	chnWebBuybackSystems :=
		chanresult.NewChanResult[map[b.SystemId]b.WebBuybackSystem](ctx, 1, 0)
	go transceiveFetchBucketData(
		ctx,
		chnWebBuybackSystems.ToSend(),
		bucketClient.ReadWebBuybackSystems,
	)

	chnWebShopLocations :=
		chanresult.NewChanResult[map[b.LocationId]b.WebShopLocation](ctx, 1, 0)
	go transceiveFetchBucketData(
		ctx,
		chnWebShopLocations.ToSend(),
		bucketClient.ReadWebShopLocations,
	)

	chnWebMarkets :=
		chanresult.NewChanResult[map[b.MarketName]b.WebMarket](ctx, 1, 0)
	go transceiveFetchBucketData(
		ctx,
		chnWebMarkets.ToSend(),
		bucketClient.ReadWebMarkets,
	)

	chnWebHRTypeMapsBuilder :=
		chanresult.NewChanResult[map[b.TypeId]b.WebHaulRouteTypeBundle](
			ctx, 1, 0,
		)
	go transceiveFetchBucketData(
		ctx,
		chnWebHRTypeMapsBuilder.ToSend(),
		bucketClient.ReadWebHaulRouteTypeMapsBuilder,
	)

	chnWebHaulRoutes :=
		chanresult.NewChanResult[map[b.WebHaulRouteSystemsKey]b.WebHaulRoute](
			ctx, 1, 0,
		)
	go transceiveFetchBucketData(
		ctx,
		chnWebHaulRoutes.ToSend(),
		bucketClient.ReadWebHaulRoutes,
	)

	if bucketData, err := chanresult.RecvOneOfEach7(
		chnWebBSTypeMapsBuilder.ToRecv(),
		chnWebSLTypeMapsBuilder.ToRecv(),
		chnWebBuybackSystems.ToRecv(),
		chnWebShopLocations.ToRecv(),
		chnWebMarkets.ToRecv(),
		chnWebHRTypeMapsBuilder.ToRecv(),
		chnWebHaulRoutes.ToRecv(),
	); err != nil {
		return b.WebBucketData{}, err
	} else {
		return b.WebBucketData{
			BuybackSystemTypeMapsBuilder: bucketData.T1,
			ShopLocationTypeMapsBuilder:  bucketData.T2,
			BuybackSystems:               bucketData.T3,
			ShopLocations:                bucketData.T4,
			Markets:                      bucketData.T5,
			HaulRouteTypeMapsBuilder:     bucketData.T6,
			HaulRoutes:                   bucketData.T7,
		}, nil
	}
}

func transceiveFetchBucketData[BD any](
	ctx context.Context,
	chnSendBucketData chanresult.ChanSendResult[BD],
	fetch func(ctx context.Context, capacity int) (BD, error),
) error {
	if bucketData, err := fetch(ctx, 0); err != nil {
		return chnSendBucketData.SendErr(err)
	} else {
		return chnSendBucketData.SendOk(bucketData)
	}
}
