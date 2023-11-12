package items

type Item interface {
	GetTypeId() int32
	GetQuantity() int64
}

func AddToMap[I Item](
	m map[int32]int64,
	items ...I,
) {
	for _, item := range items {
		m[item.GetTypeId()] += item.GetQuantity()
	}
}
