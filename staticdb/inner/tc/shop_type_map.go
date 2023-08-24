package tc

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

var KVReaderShopLocationTypeMap loader.
	LoadOnceKVReaderGobFSSlice[map[int32]int]

func InitKVReaderShopLocationTypeMap(
	chn chan<- error,
	path string,
	capacity int,
) {
	KVReaderShopLocationTypeMap = loader.
		NewLoadOnceKVReaderGobFSSlice[map[int32]int](
		path,
		capacity,
	)
	go KVReaderShopLocationTypeMap.LoadSendErr(chn)
}

// type ShopTypeData int // PricingIndex
