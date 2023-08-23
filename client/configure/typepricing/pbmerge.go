package typepricing

import (
	"fmt"

	cfg "github.com/WiggidyW/weve-esi/client/configure"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/util"
)

func ConvertPBPricing[HS util.HashSet[string]](
	pbPricing *proto.TypePricing,
	markets HS,
) (cfg.TypePricing, error) {
	if pbPricing.Percentile > 100 {
		return cfg.TypePricing{}, newError(
			"percentile must be <= 100",
		)
	} else if pbPricing.Modifier > 255 {
		return cfg.TypePricing{}, newError(
			"modifier must be <= 255",
		)
	} else if !markets.Has(pbPricing.Market) {
		return cfg.TypePricing{}, newError(
			fmt.Sprintf(
				"market '%s' does not exist",
				pbPricing.Market,
			),
		)
	}
	return cfg.TypePricing{
		IsBuy:      pbPricing.IsBuy,
		Percentile: uint8(pbPricing.Percentile),
		Modifier:   uint8(pbPricing.Modifier),
		Market:     pbPricing.Market,
	}, nil
}
