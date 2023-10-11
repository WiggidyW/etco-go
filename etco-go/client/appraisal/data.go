package appraisal

type BasicItem struct {
	TypeId   int32
	Quantity int64
}

func (bi BasicItem) GetTypeId() int32 {
	return bi.TypeId
}

func (bi BasicItem) GetQuantity() int64 {
	return bi.Quantity
}
