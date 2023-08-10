package sde

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderGroups loader.LoadOnceKVReaderGobFSSlice[Group]

func InitKVReaderGroups(chn chan<- error, path string, capacity int) {
	KVReaderGroups = loader.
		NewLoadOnceKVReaderGobFSSlice[Group](path, capacity)
	go KVReaderGroups.LoadSendErr(chn)
}

type Group struct {
	Name          string // english
	CategoryIndex *int
}
