package sde

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

var KVReaderRegions loader.LoadOnceKVReaderGobFSMap[int32, string]

func InitKVReaderRegions(chn chan<- error, path string, capacity int) {
	KVReaderRegions = loader.
		NewLoadOnceKVReaderGobFSMap[int32, string](path, capacity)
	go KVReaderRegions.LoadSendErr(chn)
}

// type Region string // english
