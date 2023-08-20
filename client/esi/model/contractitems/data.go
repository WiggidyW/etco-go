package contractitems

type ContractItemsEntry struct {
	// IsIncluded  bool   `json:"is_included"`
	// IsSingleton bool   `json:"is_singleton"`
	Quantity int32 `json:"quantity"`
	// RawQuantity *int32 `json:"raw_quantity,omitempty"`
	// RecordId    int64  `json:"record_id"`
	TypeId int32 `json:"type_id"`
}
