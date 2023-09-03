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

type CfgMergeShopLocationTypeMapsBuilderParams struct {
	Updates map[int32]*proto.CfgShopLocationTypeBundle
}

type CfgMergeShopLocationTypeMapsBuilderClient struct {
	webSTypeMapsBuilderReaderClient bucket.WebShopLocationTypeMapsBuilderReaderClient
	webSTypeMapsBuilderWriterClient bucket.WebShopLocationTypeMapsBuilderWriterClient
	webMarketReaderClient           bucket.WebMarketsReaderClient
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) Fetch(
	ctx context.Context,
	params CfgMergeShopLocationTypeMapsBuilderParams,
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
		NewChanResult[map[b.TypeId]b.WebShopLocationTypeBundle](
		ctx, 1, 0,
	).Split()
	go msbc.transceiveFetchBuilder(ctx, chnBuilderSend)

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
	if err := mergeShopLocationTypeMapsBuilder(
		builder,
		params.Updates,
		marketHashSet,
	); err != nil {
		return &CfgMergeResponse{
			// Modified: false,
			MergeError: err,
		}, nil
	}

	if err := msbc.fetchWriteUpdated(ctx, builder); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) fetchWriteUpdated(
	ctx context.Context,
	updated map[b.TypeId]b.WebShopLocationTypeBundle,
) error {
	_, err := msbc.webSTypeMapsBuilderWriterClient.Fetch(
		ctx,
		bucket.WebShopLocationTypeMapsBuilderWriterParams{
			WebShopLocationTypeMapsBuilder: updated,
		},
	)
	return err
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) fetchMarketsHashSet(
	ctx context.Context,
) (
	hashSet util.MapHashSet[string, b.WebMarket],
	err error,
) {
	if marketsRep, err := msbc.webMarketReaderClient.Fetch(
		ctx,
		bucket.WebMarketsReaderParams{},
	); err != nil {
		return nil, err
	} else {
		markets := marketsRep.Data()
		return util.MapHashSet[string, b.WebMarket](markets), nil
	}
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) transceiveFetchBuilder(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[map[b.TypeId]b.WebShopLocationTypeBundle],
) error {
	builder, err := msbc.fetchBuilder(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(builder)
	}
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) fetchBuilder(
	ctx context.Context,
) (
	builder map[b.TypeId]b.WebShopLocationTypeBundle,
	err error,
) {
	if builderRep, err := msbc.webSTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebShopLocationTypeMapsBuilderReaderParams{},
	); err != nil {
		return nil, err
	} else {
		return builderRep.Data(), nil
	}
}

func mergeShopLocationTypeMapsBuilder[HS util.HashSet[string]](
	original map[b.TypeId]b.WebShopLocationTypeBundle,
	updates map[int32]*proto.CfgShopLocationTypeBundle,
	markets HS,
) error {
	for typeId, pbShopLocationTypeBundle := range updates {
		if pbShopLocationTypeBundle == nil ||
			pbShopLocationTypeBundle.Inner == nil {
			delete(original, typeId)
		} else {
			shopLocationTypeBundle, ok := original[typeId]
			if !ok {
				shopLocationTypeBundle = make(
					b.WebShopLocationTypeBundle,
					len(pbShopLocationTypeBundle.Inner),
				)
				original[typeId] = shopLocationTypeBundle
			}
			if err := mergeShopLocationTypeBundle(
				typeId,
				shopLocationTypeBundle,
				pbShopLocationTypeBundle.Inner,
				markets,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeShopLocationTypeBundle[HS util.HashSet[string]](
	typeId b.TypeId,
	original b.WebShopLocationTypeBundle,
	updates map[string]*proto.CfgShopTypePricing,
	markets HS,
) error {
	for bundleKey, pbShopTypePricing := range updates {
		if pbShopTypePricing == nil || pbShopTypePricing.Inner == nil {
			delete(original, bundleKey)
		} else {
			shopTypePricing, err := pBToWebShopTypePricing(
				typeId,
				bundleKey,
				pbShopTypePricing.Inner,
				markets,
			)
			if err != nil {
				return err
			} else {
				original[bundleKey] = *shopTypePricing
			}
		}
	}
	return nil
}

func newPBtoWebShopTypePricingError(
	typeId int32,
	typeMapKey string,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrShopTypeInvalid{
			Err: fmt.Errorf(
				"'%d - %s': %s",
				typeId,
				typeMapKey,
				errStr,
			),
		},
	}
}

func pBToWebShopTypePricing[HS util.HashSet[string]](
	typeId b.TypeId,
	bundleKey b.BundleKey,
	pbShopTypePricing *proto.CfgTypePricing,
	markets HS,
) (
	webShopTypePricing *b.WebShopTypePricing,
	err error,
) {
	webShopTypePricing, err = PBtoWebTypePricing(
		pbShopTypePricing,
		markets,
	)

	if err != nil {
		return nil, newPBtoWebShopTypePricingError(
			typeId,
			bundleKey,
			err.Error(),
		)
	}

	return webShopTypePricing, nil
}
