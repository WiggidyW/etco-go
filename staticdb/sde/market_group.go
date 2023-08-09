package sde

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderMarketGroups loader.LoadOnceKVReaderGobFSSlice[MarketGroup]

func InitKVReaderMarketGroups(chn chan<- error, path string, capacity int) {
	kVReaderMarketGroups = loader.
		NewLoadOnceKVReaderGobFSSlice[MarketGroup](path, capacity)
	go kVReaderMarketGroups.LoadSendErr(chn)
}

type MarketGroup struct {
	Name        string // english
	ParentIndex *int
}
