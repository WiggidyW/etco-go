package tc

import "github.com/WiggidyW/weve-esi/staticdb/inner/loader"

var KVReaderBannedFlags loader.
	LoadOnceKVReaderGobFSSlice[map[string]struct{}]

func InitKVReaderBannedFlags(
	chn chan<- error,
	path string,
	capacity int,
) {
	KVReaderBannedFlags = loader.
		NewLoadOnceKVReaderGobFSSlice[map[string]struct{}](
		path,
		capacity,
	)
	go KVReaderBannedFlags.LoadSendErr(chn)
}

// type BannedFlags map[string]struct{} // set of banned flags
