package proto

import (
	"context"
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/util"
)

type CfgMergeBuybackSystemTypeMapsBuilderParams struct {
	Updates map[int32]*proto.CfgBuybackSystemTypeBundle
}

type CfgMergeBuybackSystemTypeMapsBuilderClient struct {
	webBTypeMapsBuilderReaderClient bucket.SC_WebBuybackSystemTypeMapsBuilderReaderClient
	webBTypeMapsBuilderWriterClient bucket.SAC_WebBuybackSystemTypeMapsBuilderWriterClient
	webMarketReaderClient           bucket.SC_WebMarketsReaderClient
}

func NewCfgMergeBuybackSystemTypeMapsBuilderClient(
	webBTypeMapsBuilderReaderClient bucket.SC_WebBuybackSystemTypeMapsBuilderReaderClient,
	webBTypeMapsBuilderWriterClient bucket.SAC_WebBuybackSystemTypeMapsBuilderWriterClient,
	webMarketReaderClient bucket.SC_WebMarketsReaderClient,
) CfgMergeBuybackSystemTypeMapsBuilderClient {
	return CfgMergeBuybackSystemTypeMapsBuilderClient{
		webBTypeMapsBuilderReaderClient,
		webBTypeMapsBuilderWriterClient,
		webMarketReaderClient,
	}
}

func (mbbc CfgMergeBuybackSystemTypeMapsBuilderClient) Fetch(
	ctx context.Context,
	params CfgMergeBuybackSystemTypeMapsBuilderParams,
) (
	rep *CfgMergeResponse,
	err error,
) {
	if params.Updates == nil || len(params.Updates) == 0 {
		return &CfgMergeResponse{
			// Modified: false,
			// MergeError: nil,
		}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the original builder in a goroutine
	chnBuilderSend, chnBuilderRecv := chanresult.
		NewChanResult[map[b.TypeId]b.WebBuybackSystemTypeBundle](
		ctx, 1, 0,
	).Split()
	go mbbc.transceiveFetchBuilder(ctx, chnBuilderSend)

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
	if err := mergeBuybackSystemTypeMapsBuilder(
		builder,
		params.Updates,
		marketHashSet,
	); err != nil {
		return &CfgMergeResponse{
			// Modified: false,
			MergeError: err,
		}, nil
	}

	if err := mbbc.fetchWriteUpdated(ctx, builder); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mbbc CfgMergeBuybackSystemTypeMapsBuilderClient) fetchWriteUpdated(
	ctx context.Context,
	updated map[b.TypeId]b.WebBuybackSystemTypeBundle,
) error {
	_, err := mbbc.webBTypeMapsBuilderWriterClient.Fetch(
		ctx,
		bucket.WebBuybackSystemTypeMapsBuilderWriterParams{
			WebBuybackSystemTypeMapsBuilder: updated,
		},
	)
	return err
}

func (mbbc CfgMergeBuybackSystemTypeMapsBuilderClient) fetchMarketsHashSet(
	ctx context.Context,
) (
	hashSet util.MapHashSet[string, b.WebMarket],
	err error,
) {
	if marketsRep, err := mbbc.webMarketReaderClient.Fetch(
		ctx,
		bucket.WebMarketsReaderParams{},
	); err != nil {
		return nil, err
	} else {
		markets := marketsRep.Data()
		return util.MapHashSet[string, b.WebMarket](markets), nil
	}
}

func (mbbc CfgMergeBuybackSystemTypeMapsBuilderClient) transceiveFetchBuilder(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[map[b.TypeId]b.WebBuybackSystemTypeBundle],
) error {
	builder, err := mbbc.fetchBuilder(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(builder)
	}
}

func (mbbc CfgMergeBuybackSystemTypeMapsBuilderClient) fetchBuilder(
	ctx context.Context,
) (
	builder map[b.TypeId]b.WebBuybackSystemTypeBundle,
	err error,
) {
	if builderRep, err := mbbc.webBTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebBuybackSystemTypeMapsBuilderReaderParams{},
	); err != nil {
		return nil, err
	} else {
		return builderRep.Data(), nil
	}
}

func mergeBuybackSystemTypeMapsBuilder[HS util.HashSet[string]](
	original map[b.TypeId]b.WebBuybackSystemTypeBundle,
	updates map[int32]*proto.CfgBuybackSystemTypeBundle,
	markets HS,
) error {
	for typeId, pbBuybackSystemTypeBundle := range updates {
		if pbBuybackSystemTypeBundle == nil ||
			pbBuybackSystemTypeBundle.Inner == nil ||
			len(pbBuybackSystemTypeBundle.Inner) == 0 {
			delete(original, typeId)
		} else {
			buybackSystemTypeBundle, ok := original[typeId]
			if !ok {
				buybackSystemTypeBundle = make(
					b.WebBuybackSystemTypeBundle,
					len(pbBuybackSystemTypeBundle.Inner),
				)
				original[typeId] = buybackSystemTypeBundle
			}
			if err := mergeBuybackSystemTypeBundle(
				typeId,
				buybackSystemTypeBundle,
				pbBuybackSystemTypeBundle.Inner,
				markets,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeBuybackSystemTypeBundle[HS util.HashSet[string]](
	typeId b.TypeId,
	original b.WebBuybackSystemTypeBundle,
	updates map[string]*proto.CfgBuybackTypePricing,
	markets HS,
) error {
	for bundleKey, pbBuybackTypePricing := range updates {
		if pbBuybackTypePricing == nil || (pbBuybackTypePricing.Pricing == nil &&
			pbBuybackTypePricing.ReprocessingEfficiency == 0) {
			delete(original, bundleKey)
		} else {
			buybackTypePricing, err := pBToWebBuybackTypePricing(
				typeId,
				bundleKey,
				pbBuybackTypePricing,
				markets,
			)
			if err != nil {
				return err
			} else {
				original[bundleKey] = *buybackTypePricing
			}
		}
	}
	return nil
}

func newPBtoWebBuybackTypePricingError(
	typeId int32,
	typeMapKey string,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrBuybackTypeInvalid{
			Err: fmt.Errorf(
				"'%d - %s': %s",
				typeId,
				typeMapKey,
				errStr,
			),
		},
	}
}

func pBToWebBuybackTypePricing[HS util.HashSet[string]](
	typeId b.TypeId,
	bundleKey b.BundleKey,
	pbBuybackTypePricing *proto.CfgBuybackTypePricing,
	markets HS,
) (
	webBuybackTypePricing *b.WebBuybackTypePricing,
	err error,
) {
	webBuybackTypePricing = new(b.WebBuybackTypePricing)

	var hasPricingOrReprEff bool

	if pbBuybackTypePricing.Pricing != nil {
		hasPricingOrReprEff = true

		// validate and convert .Pricing
		webBuybackTypePricing.Pricing, err = PBtoWebTypePricing(
			pbBuybackTypePricing.Pricing,
			markets,
		)

		if err != nil {
			return nil, newPBtoWebBuybackTypePricingError(
				typeId,
				bundleKey,
				err.Error(),
			)
		}
	}

	if pbBuybackTypePricing.ReprocessingEfficiency != 0 {
		hasPricingOrReprEff = true

		if pbBuybackTypePricing.ReprocessingEfficiency > 100 {
			return nil, newPBtoWebBuybackTypePricingError(
				typeId,
				bundleKey,
				"reprocessing efficiency must be <= 100",
			)
		}

		webBuybackTypePricing.ReprocessingEfficiency = uint8(
			pbBuybackTypePricing.ReprocessingEfficiency,
		)
	}

	// check our own type as a proxy for checking the actual pb variable
	if !hasPricingOrReprEff {
		return nil, newPBtoWebBuybackTypePricingError(
			typeId,
			bundleKey,
			"one of pricing or reprocessing efficiency must be set",
		)
	}

	return webBuybackTypePricing, nil
}
