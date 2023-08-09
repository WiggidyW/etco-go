package sde

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var KVReaderRegions loader.LoadOnceKVReaderGobFSMap[int32, string]

func InitKVReaderRegions(chn chan<- error, path string, capacity int) {
	KVReaderRegions = loader.
		NewLoadOnceKVReaderGobFSMap[int32, string](path, capacity)
	go KVReaderRegions.LoadSendErr(chn)
}

// type Region string // english
