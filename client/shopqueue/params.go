package shopqueue

type ShopQueueParams struct {
	// If true, will block until the DB source matches the returned data
	// returns an error instead of logging if the operation fails
	BlockOnModify bool
}

func NewShopQueueParams(blockOnModify bool) ShopQueueParams {
	return ShopQueueParams{
		BlockOnModify: blockOnModify,
	}
}
