package protoutil

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/proto"
)

func NewPBCfgBTypeMapsBuilder(
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
			NewPBCfgBuybackSystemTypeBundle(webTypeBundle)
	}
	return pbBTypeMapsBuilder
}

func NewPBCfgBuybackSystemTypeBundle(
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
			NewPBCfgBuybackTypePricing(webBTypePricing)
	}
	return pbBuybackSystemTypeBundle
}

func NewPBCfgBuybackTypePricing(
	webBuybackTypePricing b.WebBuybackTypePricing,
) (
	pbBuybackTypePricing *proto.CfgBuybackTypePricing,
) {
	if webBuybackTypePricing.Pricing != nil {
		return &proto.CfgBuybackTypePricing{
			Pricing: NewPBCfgTypePricing(
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

func NewPBCfgSTypeMapsBuilder(
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
			NewPBCfgShopLocationTypeBundle(webTypeBundle)
	}
	return pbSTypeMapsBuilder
}

func NewPBCfgShopLocationTypeBundle(
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
				Inner: NewPBCfgTypePricing(webSTypePricing),
			}
	}
	return pbShopLocationTypeBundle
}

func NewPBCfgTypePricing(
	webTypePricing b.WebTypePricing,
) (
	pbTypePricing *proto.CfgTypePricing,
) {
	return &proto.CfgTypePricing{
		IsBuy:      webTypePricing.IsBuy,
		Percentile: uint32(webTypePricing.Percentile),
		Modifier:   uint32(webTypePricing.Modifier),
		Market:     webTypePricing.MarketName,
	}
}
