package sde

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderSystems loader.LoadOnceKVReaderGobFSMap[int32, System]

func InitKVReaderSystems(chn chan<- error, path string, capacity int) {
	KVReaderSystems = loader.
		NewLoadOnceKVReaderGobFSMap[int32, System](path, capacity)
	go KVReaderSystems.LoadSendErr(chn)
}

type System struct {
	RegionId int32
	Name     string // english
}
