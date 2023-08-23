package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderBuybackSystems loader.
	LoadOnceKVReaderGobFSMap[int32, BuybackSystem]

func InitKVReaderBuybackSystems(
	chn chan<- error,
	path string,
	capacity int,
) {
	KVReaderBuybackSystems = loader.
		NewLoadOnceKVReaderGobFSMap[int32, BuybackSystem](
		path,
		capacity,
	)
	go KVReaderBuybackSystems.LoadSendErr(chn)
}

type BuybackSystem struct {
	M3Fee        *float64
	TypeMapIndex int
}
