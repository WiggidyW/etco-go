package tc

import "github.com/WiggidyW/eve-trading-co-go/staticdb/inner/loader"

var KVReaderMarket loader.
	LoadOnceKVReaderGobFSSlice[Market]

func InitKVReaderMarket(
	chn chan<- error,
	path string,
	capacity int,
) {
	KVReaderMarket = loader.
		NewLoadOnceKVReaderGobFSSlice[Market](
		path,
		capacity,
	)
	go KVReaderMarket.LoadSendErr(chn)
}

type Market struct {
	Name         string // user-defined
	RefreshToken *string
	LocationId   int64
	IsStructure  bool
}
