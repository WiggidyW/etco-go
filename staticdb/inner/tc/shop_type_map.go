package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderShopTypeMap loader.LoadOnceKVReaderGobFSSlice[map[int32]int]

func InitKVReaderShopTypeMap(chn chan<- error, path string, capacity int) {
	KVReaderShopTypeMap = loader.
		NewLoadOnceKVReaderGobFSSlice[map[int32]int](path, capacity)
	go KVReaderShopTypeMap.LoadSendErr(chn)
}

// type ShopTypeData int // PricingIndex
