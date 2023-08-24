package tc

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

var KVReaderShopLocations loader.
	LoadOnceKVReaderGobFSMap[int64, ShopLocation]

func InitKVReaderShopLocations(
	chn chan<- error,
	path string,
	capacity int,
) {
	KVReaderShopLocations = loader.
		NewLoadOnceKVReaderGobFSMap[int64, ShopLocation](
		path,
		capacity,
	)
	go KVReaderShopLocations.LoadSendErr(chn)
}

type ShopLocation struct {
	BannedFlagsIndex *int
	TypeMapIndex     int
}
