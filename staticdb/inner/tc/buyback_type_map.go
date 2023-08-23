package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderBuybackSystemTypeMap loader.
	LoadOnceKVReaderGobFSSlice[map[int32]BuybackTypeData]

func InitKVReaderBuybackSystemTypeMap(
	chn chan<- error,
	path string,
	capacity int,
) {
	KVReaderBuybackSystemTypeMap = loader.
		NewLoadOnceKVReaderGobFSSlice[map[int32]BuybackTypeData](
		path,
		capacity,
	)
	go KVReaderBuybackSystemTypeMap.LoadSendErr(chn)
}

type BuybackTypeData struct {
	ReprocessingEfficiency uint8 // 0 = nil, 1 - 100 = efficiency
	PricingIndex           *int
}
