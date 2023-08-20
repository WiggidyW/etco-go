package rawshopqueue

const (
	CACHE_KEY     = "shopqueue"
	COLLECTION_ID = "shop_queue"
	DOC_ID        = "shop_queue"
	FIELD_ID      = "ShopQueue"
)

type ShopQueue struct {
	ShopQueue []string `firestore:"shop_queue"`
}
