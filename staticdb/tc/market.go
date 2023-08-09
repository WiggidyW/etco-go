package tc

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderMarket loader.LoadOnceKVReaderGobFSSlice[Market]

func InitKVReaderMarket(chn chan<- error, path string, capacity int) {
	kVReaderMarket = loader.
		NewLoadOnceKVReaderGobFSSlice[Market](path, capacity)
	go kVReaderMarket.LoadSendErr(chn)
}

type Market struct {
	Name         string // user-defined
	RefreshToken *string
	LocationId   int64
	IsStructure  bool
}
