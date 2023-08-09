package tc

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderBuybackSystem loader.LoadOnceKVReaderGobFSMap[int32, BuybackSystem]

func InitKVReaderBuybackSystem(chn chan<- error, path string, capacity int) {
	kVReaderBuybackSystem = loader.
		NewLoadOnceKVReaderGobFSMap[int32, BuybackSystem](path, capacity)
	go kVReaderBuybackSystem.LoadSendErr(chn)
}

type BuybackSystem struct {
	M3Fee        float64
	TypeMapIndex int
}
