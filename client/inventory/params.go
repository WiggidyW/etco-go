package inventory

type InventoryParams struct {
	LocationId int64
}

func NewInventoryParams(
	locationId int64,
) InventoryParams {
	return InventoryParams{
		LocationId: locationId,
	}
}
