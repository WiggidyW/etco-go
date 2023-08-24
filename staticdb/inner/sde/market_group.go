package sde

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

var KVReaderMarketGroups loader.LoadOnceKVReaderGobFSSlice[MarketGroup]

func InitKVReaderMarketGroups(chn chan<- error, path string, capacity int) {
	KVReaderMarketGroups = loader.
		NewLoadOnceKVReaderGobFSSlice[MarketGroup](path, capacity)
	go KVReaderMarketGroups.LoadSendErr(chn)
}

type MarketGroup struct {
	Name        string // english
	NumParents  uint8
	ParentIndex *int
}
