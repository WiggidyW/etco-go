package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderBuybackTypeMap loader.LoadOnceKVReaderGobFSSlice[map[int32]BuybackTypeData]

func InitKVReaderBuybackTypeMap(chn chan<- error, path string, capacity int) {
	KVReaderBuybackTypeMap = loader.
		NewLoadOnceKVReaderGobFSSlice[map[int32]BuybackTypeData](path, capacity)
	go KVReaderBuybackTypeMap.LoadSendErr(chn)
}

type BuybackTypeData struct {
	ReprocessingEfficiency uint8 // 0 = nil, 1 - 100 = efficiency
	PricingIndex           *int
}
