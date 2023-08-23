package pbmerge

import (
	cfg "github.com/WiggidyW/weve-esi/client/configure"
	"github.com/WiggidyW/weve-esi/client/configure/typepricing"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/util"
)

func ConvertPBShopTypePricing[HS util.HashSet[string]](
	typeId cfg.TypeId,
	bundleKey cfg.BundleKey,
	pbShopTypePricing *proto.TypePricing,
	markets HS,
) (cfg.ShopTypePricing, error) {
	if typePricing, err := typepricing.ConvertPBPricing(
		pbShopTypePricing,
		markets,
	); err != nil {
		return cfg.ShopTypePricing{}, newError(
			typeId,
			bundleKey,
			err.Error(),
		)
	} else {
		return typePricing, nil
	}
}

func MergeShopLocationTypeBundle[HS util.HashSet[string]](
	typeId cfg.TypeId,
	original cfg.ShopLocationTypeBundle,
	updates map[string]*proto.ShopTypePricing,
	markets HS,
) error {
	for bundleKey, pbShopTypePricing := range updates {
		if pbShopTypePricing == nil || pbShopTypePricing.Inner == nil {
			delete(original, bundleKey)
		} else {
			shopTypePricing, err := ConvertPBShopTypePricing(
				typeId,
				bundleKey,
				pbShopTypePricing.Inner,
				markets,
			)
			if err != nil {
				return err
			} else {
				original[bundleKey] = shopTypePricing
			}
		}
	}
	return nil
}

func MergeShopLocationTypeMapsBuilder[HS util.HashSet[string]](
	original cfg.ShopLocationTypeMapsBuilder,
	updates map[int32]*proto.ShopLocationTypeBundle,
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
					cfg.ShopLocationTypeBundle,
					len(pbShopLocationTypeBundle.Inner),
				)
				original[typeId] = shopLocationTypeBundle
			}
			if err := MergeShopLocationTypeBundle(
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
