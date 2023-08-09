package sde

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderTypeData loader.LoadOnceKVReaderGobFSMap[int32, TypeData]

func InitKVReaderTypeData(chn chan<- error, path string, capacity int) {
	kVReaderTypeData = loader.
		NewLoadOnceKVReaderGobFSMap[int32, TypeData](path, capacity)
	go kVReaderTypeData.LoadSendErr(chn)
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
