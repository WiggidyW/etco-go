package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var kVReaderTypeData loader.LoadOnceKVReaderGobFSMap[int32, TypeData]

func InitKVReaderTypeData(chn chan<- error, path string, capacity int) {
	kVReaderTypeData = loader.
		NewLoadOnceKVReaderGobFSMap[int32, TypeData](path, capacity)
	go kVReaderTypeData.LoadSendErr(chn)
}

type TypeData struct {
	BuybackReprocessingEfficiency uint8 // 0 = nil, 1 - 100 = efficiency
	BuybackPricingIndex           *int
	ShopPricingIndex              *int
}
