package sde

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderCategories loader.LoadOnceKVReaderGobFSSlice[string]

func InitKVReaderCategories(chn chan<- error, path string, capacity int) {
	KVReaderCategories = loader.
		NewLoadOnceKVReaderGobFSSlice[string](path, capacity)
	go KVReaderCategories.LoadSendErr(chn)
}

// type Category string // english
