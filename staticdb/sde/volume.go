package sde

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderVolumes loader.LoadOnceKVReaderGobFSSlice[float64]

func InitKVReaderVolumes(chn chan<- error, path string, capacity int) {
	kVReaderVolumes = loader.
		NewLoadOnceKVReaderGobFSSlice[float64](path, capacity)
	go kVReaderVolumes.LoadSendErr(chn)
}

// type Category string // english
