package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func GetSystemInfo(systemId b.SystemId) (systemInfo *b.System) {
	v, exists := kvreader_.KVReaderSystems.Get(systemId)
	if exists {
		return &v
	} else {
		return nil
	}
}
