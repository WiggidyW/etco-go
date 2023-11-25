package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func UNSAFE_GetSystemInfos() map[b.SystemId]b.System {
	return kvreader_.KVReaderSystems.UnsafeGetInner().UnsafeGetInner()
}

func GetSystemIds() []b.SystemId {
	systemsMap := UNSAFE_GetSystemInfos()
	systemIds := make([]int32, len(systemsMap))
	i := 0
	for systemId := range systemsMap {
		systemIds[i] = systemId
		i++
	}
	return systemIds
}

func GetSystemInfo(systemId b.SystemId) (systemInfo *b.System) {
	v, exists := kvreader_.KVReaderSystems.Get(systemId)
	if exists {
		return &v
	} else {
		return nil
	}
}
