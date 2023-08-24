package sde

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

// only includes types with published = true and marketGroupID != null
var KVReaderTypeIDs loader.LoadOnceKVReaderGobFSMap[string, int32]

func InitKVReaderTypeIDs(chn chan<- error, path string, capacity int) {
	KVReaderTypeIDs = loader.
		NewLoadOnceKVReaderGobFSMap[string, int32](path, capacity)
	go KVReaderTypeIDs.LoadSendErr(chn)
}

// type TypeID int32
