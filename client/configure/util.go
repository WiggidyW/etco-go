package configure

type Validatable[K any] interface {
	Validate(k K) error
}

func MergeMaps[K comparable, V Validatable[K]](
	original map[K]V,
	updates map[K]V,
) (map[K]V, error) {
	for k, v := range updates {
		if err := v.Validate(k); err != nil {
			return nil, err
		}
		original[k] = v
	}
	return original, nil
}
