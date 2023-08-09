package sde

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderCategories loader.LoadOnceKVReaderGobFSSlice[string]

func InitKVReaderCategories(chn chan<- error, path string, capacity int) {
	kVReaderCategories = loader.
		NewLoadOnceKVReaderGobFSSlice[string](path, capacity)
	go kVReaderCategories.LoadSendErr(chn)
}

// type Category string // english
