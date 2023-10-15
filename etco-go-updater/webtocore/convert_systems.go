package webtocore

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

func convertWebBuybackSystems(
	webBuybackSystems map[b.SystemId]b.WebBuybackSystem,
	coreBSTypeMapsIndexMap map[b.BundleKey]int,
) (
	coreBuybackSystems map[b.SystemId]b.BuybackSystem,
	err error,
) {
	coreBuybackSystems = make(
		map[b.SystemId]b.BuybackSystem,
		len(webBuybackSystems),
	)

	for systemId, webBuybackSystem := range webBuybackSystems {
		typeMapIndex, ok := coreBSTypeMapsIndexMap[webBuybackSystem.
			BundleKey]
		if !ok {
			return nil, fmt.Errorf(
				"BuybackSystem %d has invalid BundleKey %s",
				systemId,
				webBuybackSystem.BundleKey,
			)
		}
		coreBuybackSystems[systemId] = b.BuybackSystem{
			M3Fee:        webBuybackSystem.M3Fee,
			TaxRate:      webBuybackSystem.TaxRate,
			TypeMapIndex: typeMapIndex,
		}
	}

	return coreBuybackSystems, nil
}
