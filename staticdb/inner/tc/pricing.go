package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderPricing loader.LoadOnceKVReaderGobFSSlice[Pricing]

func InitKVReaderPricing(chn chan<- error, path string, capacity int) {
	KVReaderPricing = loader.
		NewLoadOnceKVReaderGobFSSlice[Pricing](path, capacity)
	go KVReaderPricing.LoadSendErr(chn)
}

type Pricing struct {
	IsBuy       bool
	Percentile  uint8 // 0 - 100
	Modifier    uint8 // 1 - 255
	MarketIndex int
}
