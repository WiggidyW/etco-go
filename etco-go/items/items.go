package items

type IBasicItem interface {
	GetTypeId() int32
	GetQuantity() int64
}

type BasicItem struct {
	TypeId   int32
	Quantity int64
}

func (bi BasicItem) GetTypeId() int32   { return bi.TypeId }
func (bi BasicItem) GetQuantity() int64 { return bi.Quantity }

func AddToMap[I IBasicItem](
	m map[int32]int64,
	items ...I,
) {
	for _, item := range items {
		m[item.GetTypeId()] += item.GetQuantity()
	}
}

func NewBasicItems[I IBasicItem](items []I) []BasicItem {
	basicItems := make([]BasicItem, len(items))
	for i, item := range items {
		basicItems[i] = BasicItem{
			TypeId:   item.GetTypeId(),
			Quantity: item.GetQuantity(),
		}
	}
	return basicItems
}
