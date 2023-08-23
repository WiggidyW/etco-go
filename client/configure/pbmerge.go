package configure

func ExtractBuilderBundleKeys[V any](
	builder map[int32]map[string]V,
) map[string]struct{} {
	bundleKeys := make(map[string]struct{})
	for _, bundle := range builder {
		for bundleKey := range bundle {
			bundleKeys[bundleKey] = struct{}{}
		}
	}
	return bundleKeys
}

// func ExtractBuybackBuilderBundleKeys(builder BuybackSystemTypeMapsBuilder) map[string]struct{} {
// 	return ExtractBuilderBundleKeys(builder)
// }
