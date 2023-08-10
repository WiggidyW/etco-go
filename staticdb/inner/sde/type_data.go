package sde

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderTypeData loader.LoadOnceKVReaderGobFSMap[int32, TypeData]

func InitKVReaderTypeData(chn chan<- error, path string, capacity int) {
	KVReaderTypeData = loader.
		NewLoadOnceKVReaderGobFSMap[int32, TypeData](path, capacity)
	go KVReaderTypeData.LoadSendErr(chn)
}

type TypeData struct {
	ReprocessedMaterials []ReprocessedMaterial
	Name                 string // english
	GroupIndex           *int
	MarketGroupIndex     *int
	VolumeIndex          *int
}

type ReprocessedMaterial struct {
	TypeId   int32
	Quantity float64
}
