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
	WEB_S_TYPEMAPSBUILDER_BUF_CAP          int           = 0
	WEB_S_TYPEMAPSBUILDER_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebShopLocationTypeMapsBuilder = cache.RegisterType[map[b.TypeId]b.WebShopLocationTypeBundle]("webshoplocationtypemapsbuilder", WEB_S_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebShopLocationTypeMapsBuilder(
	x cache.Context,
) (
	rep map[b.TypeId]b.WebShopLocationTypeBundle,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopLocationTypeMapsBuilder,
		WEB_S_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}

func GetWebShopLocationActiveMarkets(
	x cache.Context,
) (
	rep map[string]struct{},
	expires time.Time,
	err error,
) {
	var webShopLocationTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle
	webShopLocationTypeMapsBuilder, expires, err =
		GetWebShopLocationTypeMapsBuilder(x)
	if err == nil {
		rep = extractShopBuilderActiveMarkets(
			webShopLocationTypeMapsBuilder,
		)
	}
	return rep, expires, err
}

func ProtoGetWebShopLocationTypeMapsBuilder(
	x cache.Context,
) (
	rep map[int32]*proto.CfgShopLocationTypeBundle,
	expires time.Time,
	err error,
) {
	var webShopLocationTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle
	webShopLocationTypeMapsBuilder, expires, err =
		GetWebShopLocationTypeMapsBuilder(x)
	if err == nil {
		rep = WebShopLocationTypeMapsBuilderToProto(
			webShopLocationTypeMapsBuilder,
		)
	}
	return rep, expires, err
}

func SetWebShopLocationTypeMapsBuilder(
	x cache.Context,
	rep map[b.TypeId]b.WebShopLocationTypeBundle,
) (
	err error,
) {
	lock := cacheprefetch.ServerLock(
		keys.CacheKeyWebShopBundleKeys,
		keys.TypeStrWebShopBundleKeys,
	)
	return set(
		x,
		client.WriteWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopLocationTypeMapsBuilder,
		WEB_S_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
		&lock,
	)
}

func ProtoMergeSetWebShopLocationTypeMapsBuilder(
	x cache.Context,
	updates map[int32]*proto.CfgShopLocationTypeBundle,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}
	return protoMergeSetTypeMapsBuilder(
		x,
		updates,
		GetWebShopLocationTypeMapsBuilder,
		ProtoMergeShopLocationTypeMapsBuilder,
		SetWebShopLocationTypeMapsBuilder,
	)
}

// // Active Markets

func extractShopBuilderActiveMarkets(
	sbuilder map[b.TypeId]b.WebShopLocationTypeBundle,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range sbuilder {
		for _, sTypePricing := range bundle {
			activeMarkets[sTypePricing.MarketName] = struct{}{}
		}
	}
	return activeMarkets
}

// // To Proto

func WebShopLocationTypeMapsBuilderToProto(
	webBTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle,
) (
	pbBTypeMapsBuilder map[int32]*proto.CfgShopLocationTypeBundle,
) {
	return newPBCfgSTypeMapsBuilder(webBTypeMapsBuilder)
}

func newPBCfgSTypeMapsBuilder(
	webSTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle,
) (
	pbSTypeMapsBuilder map[int32]*proto.CfgShopLocationTypeBundle,
) {
	pbSTypeMapsBuilder = make(
		map[int32]*proto.CfgShopLocationTypeBundle,
		len(webSTypeMapsBuilder),
	)
	for typeId, webTypeBundle := range webSTypeMapsBuilder {
		pbSTypeMapsBuilder[typeId] =
			newPBCfgShopLocationTypeBundle(webTypeBundle)
	}
	return pbSTypeMapsBuilder
}

func newPBCfgShopLocationTypeBundle(
	webShopLocationTypeBundle b.WebShopLocationTypeBundle,
) (
	pbShopLocationTypeBundle *proto.CfgShopLocationTypeBundle,
) {
	pbShopLocationTypeBundle = &proto.CfgShopLocationTypeBundle{
		Inner: make(
			map[string]*proto.CfgShopTypePricing,
			len(webShopLocationTypeBundle),
		),
	}
	for bundleKey, webSTypePricing := range webShopLocationTypeBundle {
		pbShopLocationTypeBundle.Inner[bundleKey] =
			&proto.CfgShopTypePricing{
				Inner: newPBCfgTypePricing(webSTypePricing),
			}
	}
	return pbShopLocationTypeBundle
}

// // Merge

func ProtoMergeShopLocationTypeMapsBuilder[HSV any](
	original map[b.TypeId]b.WebShopLocationTypeBundle,
	updates map[int32]*proto.CfgShopLocationTypeBundle,
	markets map[string]HSV,
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

func mergeShopLocationTypeBundle[HSV any](
	typeId b.TypeId,
	original b.WebShopLocationTypeBundle,
	updates map[string]*proto.CfgShopTypePricing,
	markets map[string]HSV,
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

func pBToWebShopTypePricing[HSV any](
	typeId b.TypeId,
	bundleKey b.BundleKey,
	pbShopTypePricing *proto.CfgTypePricing,
	markets map[string]HSV,
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
