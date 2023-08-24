package sde

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

var KVReaderVolumes loader.LoadOnceKVReaderGobFSSlice[float64]

func InitKVReaderVolumes(chn chan<- error, path string, capacity int) {
	KVReaderVolumes = loader.
		NewLoadOnceKVReaderGobFSSlice[float64](path, capacity)
	go KVReaderVolumes.LoadSendErr(chn)
}

// type Category string // english
