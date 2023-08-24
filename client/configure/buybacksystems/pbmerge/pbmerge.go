package pbmerge

import (
	"fmt"

	cfg "github.com/WiggidyW/eve-trading-co-go/client/configure"
	"github.com/WiggidyW/eve-trading-co-go/proto"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

func ConvertPBBuybackSystem(
	pbBuybackSystem *proto.BuybackSystem,
) cfg.BuybackSystem {
	return cfg.BuybackSystem{
		BundleKey: pbBuybackSystem.BundleKey,
		M3Fee:     pbBuybackSystem.M3Fee,
	}
}

func MergeBuybackSystems[HS util.HashSet[string]](
	original cfg.BuybackSystems,
	updates map[int32]*proto.BuybackSystem,
	bundleKeys HS,
) error {
	// if updates == nil || len(updates.Inner) == 0 {
	// 	return false, nil
	// }
	for systemId, pbBuybackSystem := range updates {
		if pbBuybackSystem == nil {
			delete(original, systemId)
		} else if !bundleKeys.Has(pbBuybackSystem.BundleKey) {
			return newError(systemId, fmt.Sprintf(
				"type map key '%s' does not exist",
				pbBuybackSystem.BundleKey,
			))
		} else {
			original[systemId] = ConvertPBBuybackSystem(
				pbBuybackSystem,
			)
		}
	}
	return nil
}
