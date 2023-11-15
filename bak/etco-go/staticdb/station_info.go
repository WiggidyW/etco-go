package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func GetStationInfo(stationId b.StationId) (stationInfo *b.Station) {
	v, exists := kvreader_.KVReaderStations.Get(stationId)
	if exists {
		return &v
	} else {
		return nil
	}
}
