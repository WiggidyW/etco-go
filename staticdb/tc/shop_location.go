package tc

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderShopLocation loader.LoadOnceKVReaderGobFSMap[int64, ShopLocation]

func InitKVReaderShopLocation(chn chan<- error, path string, capacity int) {
	kVReaderShopLocation = loader.
		NewLoadOnceKVReaderGobFSMap[int64, ShopLocation](path, capacity)
	go kVReaderShopLocation.LoadSendErr(chn)
}

type ShopLocation struct {
	BannedFlagsIndex *int
	TypeMapIndex     int
}
