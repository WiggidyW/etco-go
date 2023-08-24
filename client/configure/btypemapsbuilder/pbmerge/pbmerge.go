package pbmerge

import (
	cfg "github.com/WiggidyW/eve-trading-co-go/client/configure"
	"github.com/WiggidyW/eve-trading-co-go/client/configure/typepricing"
	"github.com/WiggidyW/eve-trading-co-go/proto"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

func convertPBBuybackTypePricing[HS util.HashSet[string]](
	typeId cfg.TypeId,
	bundleKey cfg.BundleKey,
	pbBuybackTypePricing *proto.BuybackTypePricing,
	markets HS,
) (cfg.BuybackTypePricing, error) {
	buybackTypePricing := new(cfg.BuybackTypePricing)

	if pbBuybackTypePricing.Pricing != nil {
		// validate and convert .Pricing
		if typePricing, err := typepricing.ConvertPBPricing(
			pbBuybackTypePricing.Pricing,
			markets,
		); err != nil {
			return cfg.BuybackTypePricing{}, newError(
				typeId,
				bundleKey,
				err.Error(),
			)
		} else {
			buybackTypePricing.Pricing = &typePricing
		}
	}

	if pbBuybackTypePricing.ReprocessingEfficiency != 0 {
		if pbBuybackTypePricing.ReprocessingEfficiency > 100 {
			return cfg.BuybackTypePricing{}, newError(
				typeId,
				bundleKey,
				"reprocessing efficiency must be <= 100",
			)
		} else {
			buybackTypePricing.ReprEff = uint8(
				pbBuybackTypePricing.ReprocessingEfficiency,
			)
		}
	}

	// check our own type as a proxy for checking the actual pb variable
	if buybackTypePricing.Pricing == nil &&
		buybackTypePricing.ReprEff == 0 {
		return cfg.BuybackTypePricing{}, newError(
			typeId,
			bundleKey,
			"one of pricing or reprocessing efficiency must be set",
		)
	}

	return *buybackTypePricing, nil
}

func mergeBuybackSystemTypeBundle[HS util.HashSet[string]](
	typeId cfg.TypeId,
	original cfg.BuybackSystemTypeBundle,
	updates map[string]*proto.BuybackTypePricing,
	markets HS,
) error {
	for bundleKey, pbBuybackTypePricing := range updates {
		if pbBuybackTypePricing == nil {
			delete(original, bundleKey)
		} else {
			buybackTypePricing, err := convertPBBuybackTypePricing(
				typeId,
				bundleKey,
				pbBuybackTypePricing,
				markets,
			)
			if err != nil {
				return err
			} else {
				original[bundleKey] = buybackTypePricing
			}
		}
	}
	return nil
}

func MergeBuybackSystemTypeMapsBuilder[HS util.HashSet[string]](
	original cfg.BuybackSystemTypeMapsBuilder,
	updates map[int32]*proto.BuybackSystemTypeBundle,
	markets HS,
) error {
	for typeId, pbBuybackSystemTypeBundle := range updates {
		if pbBuybackSystemTypeBundle == nil ||
			pbBuybackSystemTypeBundle.Inner == nil {
			delete(original, typeId)
		} else {
			buybackSystemTypeBundle, ok := original[typeId]
			if !ok {
				buybackSystemTypeBundle = make(
					cfg.BuybackSystemTypeBundle,
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
