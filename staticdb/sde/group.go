package sde

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderGroups loader.LoadOnceKVReaderGobFSSlice[Group]

func InitKVReaderGroups(chn chan<- error, path string, capacity int) {
	kVReaderGroups = loader.
		NewLoadOnceKVReaderGobFSSlice[Group](path, capacity)
	go kVReaderGroups.LoadSendErr(chn)
}

type Group struct {
	Name          string // english
	CategoryIndex *int
}
