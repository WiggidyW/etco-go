package bucket

import (
	"fmt"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/proto"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_B_TYPEMAPSBUILDER_BUF_CAP    int           = 0
	WEB_B_TYPEMAPSBUILDER_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebBuybackSystemTypeMapsBuilder = cache.RegisterType[map[b.TypeId]b.WebBuybackSystemTypeBundle]("webbuybacksystemtypemapsbuilder", WEB_B_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebBuybackSystemTypeMapsBuilder(
	x cache.Context,
) (
	rep map[b.TypeId]b.WebBuybackSystemTypeBundle,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackSystemTypeMapsBuilder,
		WEB_B_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}

func GetWebBuybackSystemActiveMarkets(
	x cache.Context,
) (
	rep map[string]struct{},
	expires time.Time,
	err error,
) {
	var webBuybackSystemTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle
	webBuybackSystemTypeMapsBuilder, expires, err =
		GetWebBuybackSystemTypeMapsBuilder(x)
	if err == nil {
		rep = extractBuybackBuilderActiveMarkets(
			webBuybackSystemTypeMapsBuilder,
		)
	}
	return rep, expires, err
}

func ProtoGetWebBuybackSystemTypeMapsBuilder(
	x cache.Context,
) (
	rep map[int32]*proto.CfgBuybackSystemTypeBundle,
	expires time.Time,
	err error,
) {
	var webBuybackSystemTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle
	webBuybackSystemTypeMapsBuilder, expires, err =
		GetWebBuybackSystemTypeMapsBuilder(x)
	if err == nil {
		rep = WebBuybackSystemTypeMapsBuilderToProto(
			webBuybackSystemTypeMapsBuilder,
		)
	}
	return rep, expires, err
}

func SetWebBuybackSystemTypeMapsBuilder(
	x cache.Context,
	rep map[b.TypeId]b.WebBuybackSystemTypeBundle,
) (
	err error,
) {
	lock := cacheprefetch.ServerLock(
		keys.CacheKeyWebBuybackBundleKeys,
		keys.TypeStrWebBuybackBundleKeys,
	)
	return set(
		x,
		client.WriteWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackSystemTypeMapsBuilder,
		WEB_B_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
		&lock,
	)
}

func ProtoMergeSetWebBuybackSystemTypeMapsBuilder(
	x cache.Context,
	updates map[int32]*proto.CfgBuybackSystemTypeBundle,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}
	return protoMergeSetTypeMapsBuilder(
		x,
		updates,
		GetWebBuybackSystemTypeMapsBuilder,
		ProtoMergeBuybackSystemTypeMapsBuilder,
		SetWebBuybackSystemTypeMapsBuilder,
	)
}

// // Active Markets

func extractBuybackBuilderActiveMarkets(
	bbuilder map[b.TypeId]b.WebBuybackSystemTypeBundle,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range bbuilder {
		for _, bTypePricing := range bundle {
			if bTypePricing.Pricing == nil {
				continue
			}
			activeMarkets[bTypePricing.Pricing.MarketName] =
				struct{}{}
		}
	}
	return activeMarkets
}

// // ToProto

func WebBuybackSystemTypeMapsBuilderToProto(
	webBTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle,
) (
	pbBTypeMapsBuilder map[int32]*proto.CfgBuybackSystemTypeBundle,
) {
	return newPBCfgBTypeMapsBuilder(webBTypeMapsBuilder)
}

func newPBCfgBTypeMapsBuilder(
	webBTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle,
) (
	pbBTypeMapsBuilder map[int32]*proto.CfgBuybackSystemTypeBundle,
) {
	pbBTypeMapsBuilder = make(
		map[int32]*proto.CfgBuybackSystemTypeBundle,
		len(webBTypeMapsBuilder),
	)
	for typeId, webTypeBundle := range webBTypeMapsBuilder {
		pbBTypeMapsBuilder[typeId] =
			newPBCfgBuybackSystemTypeBundle(webTypeBundle)
	}
	return pbBTypeMapsBuilder
}

func newPBCfgBuybackSystemTypeBundle(
	webBuybackSystemTypeBundle b.WebBuybackSystemTypeBundle,
) (
	pbBuybackSystemTypeBundle *proto.CfgBuybackSystemTypeBundle,
) {
	pbBuybackSystemTypeBundle = &proto.CfgBuybackSystemTypeBundle{
		Inner: make(
			map[string]*proto.CfgBuybackTypePricing,
			len(webBuybackSystemTypeBundle),
		),
	}
	for bundleKey, webBTypePricing := range webBuybackSystemTypeBundle {
		pbBuybackSystemTypeBundle.Inner[bundleKey] =
			newPBCfgBuybackTypePricing(webBTypePricing)
	}
	return pbBuybackSystemTypeBundle
}

func newPBCfgBuybackTypePricing(
	webBuybackTypePricing b.WebBuybackTypePricing,
) (
	pbBuybackTypePricing *proto.CfgBuybackTypePricing,
) {
	if webBuybackTypePricing.Pricing != nil {
		return &proto.CfgBuybackTypePricing{
			Pricing: newPBCfgTypePricing(
				*webBuybackTypePricing.Pricing,
			),
			ReprocessingEfficiency: uint32(
				webBuybackTypePricing.ReprocessingEfficiency,
			),
		}
	} else {
		return &proto.CfgBuybackTypePricing{
			// Pricing: nil,
			ReprocessingEfficiency: uint32(
				webBuybackTypePricing.ReprocessingEfficiency,
			),
		}
	}
}

// // Merge

func ProtoMergeBuybackSystemTypeMapsBuilder[HSV any](
	original map[b.TypeId]b.WebBuybackSystemTypeBundle,
	updates map[int32]*proto.CfgBuybackSystemTypeBundle,
	markets map[string]HSV,
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

func mergeBuybackSystemTypeBundle[HSV any](
	typeId b.TypeId,
	original b.WebBuybackSystemTypeBundle,
	updates map[string]*proto.CfgBuybackTypePricing,
	markets map[string]HSV,
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

func pBToWebBuybackTypePricing[HSV any](
	typeId b.TypeId,
	bundleKey b.BundleKey,
	pbBuybackTypePricing *proto.CfgBuybackTypePricing,
	markets map[string]HSV,
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
