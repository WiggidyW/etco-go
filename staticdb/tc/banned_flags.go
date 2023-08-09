package tc

import "github.com/WiggidyW/weve-esi/staticdb/loader"

var kVReaderBannedFlags loader.LoadOnceKVReaderGobFSSlice[map[string]struct{}]

func InitKVReaderBannedFlags(chn chan<- error, path string, capacity int) {
	kVReaderBannedFlags = loader.
		NewLoadOnceKVReaderGobFSSlice[map[string]struct{}](path, capacity)
	go kVReaderBannedFlags.LoadSendErr(chn)
}

// type BannedFlags map[string]struct{} // set of banned flags
