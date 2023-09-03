package protoutil

type BasicItem interface {
	GetTypeId() int32
	GetQuantity() int64
}
