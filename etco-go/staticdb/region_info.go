package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func GetRegionInfo(regionId b.RegionId) (regionInfo *b.RegionName) {
	v, exists := kvreader_.KVReaderRegions.Get(regionId)
	if exists {
		return &v
	} else {
		return nil
	}
}
