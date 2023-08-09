package tc

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderShopTypeMap loader.LoadOnceKVReaderGobFSSlice[map[int32]int]

func InitkVReaderShopTypeMap(chn chan<- error, path string, capacity int) {
	kVReaderShopTypeMap = loader.
		NewLoadOnceKVReaderGobFSSlice[map[int32]int](path, capacity)
	go kVReaderShopTypeMap.LoadSendErr(chn)
}

// type ShopTypeData int // PricingIndex
