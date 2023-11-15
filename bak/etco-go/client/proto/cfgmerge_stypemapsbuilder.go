package proto

import (
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/util"
)

type CfgMergeShopLocationTypeMapsBuilderParams struct {
	Updates map[int32]*proto.CfgShopLocationTypeBundle
}

type CfgMergeShopLocationTypeMapsBuilderClient struct{}

func NewCfgMergeShopLocationTypeMapsBuilderClient() CfgMergeShopLocationTypeMapsBuilderClient {
	return CfgMergeShopLocationTypeMapsBuilderClient{}
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) Fetch(
	x cache.Context,
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

	x, cancel := x.WithCancel()
	defer cancel()

	// fetch the original builder in a goroutine
	chnBuilderSend, chnBuilderRecv := chanresult.
		NewChanResult[map[b.TypeId]b.WebShopLocationTypeBundle](
		x.Ctx(), 1, 0,
	).Split()
	go msbc.transceiveFetchBuilder(x, chnBuilderSend)

	// fetch markets (used for update validation - ensures markets exist)
	marketHashSet, err := msbc.fetchMarketsHashSet(x)
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

	if err := msbc.fetchWriteUpdated(x, builder); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) fetchWriteUpdated(
	x cache.Context,
	updated map[b.TypeId]b.WebShopLocationTypeBundle,
) error {
	return bucket.SetWebShopLocationTypeMapsBuilder(x, updated)
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) fetchMarketsHashSet(
	x cache.Context,
) (
	hashSet util.MapHashSet[string, b.WebMarket],
	err error,
) {
	markets, _, err := bucket.GetWebMarkets(x)
	if err == nil {
		hashSet = util.MapHashSet[string, b.WebMarket](markets)
	}
	return hashSet, err
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) transceiveFetchBuilder(
	x cache.Context,
	chnSend chanresult.ChanSendResult[map[b.TypeId]b.WebShopLocationTypeBundle],
) error {
	builder, err := msbc.fetchBuilder(x)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(builder)
	}
}

func (msbc CfgMergeShopLocationTypeMapsBuilderClient) fetchBuilder(
	x cache.Context,
) (
	builder map[b.TypeId]b.WebShopLocationTypeBundle,
	err error,
) {
	builder, _, err = bucket.GetWebShopLocationTypeMapsBuilder(x)
	return builder, err
}

func mergeShopLocationTypeMapsBuilder[HS util.HashSet[string]](
	original map[b.TypeId]b.WebShopLocationTypeBundle,
	updates map[int32]*proto.CfgShopLocationTypeBundle,
	markets HS,
) error {
	for typeId, pbShopLocationTypeBundle := range updates {
		if pbShopLocationTypeBundle == nil ||
			pbShopLocationTypeBundle.Inner == nil ||
			len(pbShopLocationTypeBundle.Inner) == 0 {
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
