package sde

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderStations loader.LoadOnceKVReaderGobFSMap[int32, Station]

func InitKVReaderStations(chn chan<- error, path string, capacity int) {
	KVReaderStations = loader.
		NewLoadOnceKVReaderGobFSMap[int32, Station](path, capacity)
	go KVReaderStations.LoadSendErr(chn)
}

type Station struct {
	SystemId int32
	Name     string // english
}
