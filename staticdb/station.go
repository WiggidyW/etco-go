package staticdb

import (
	"github.com/WiggidyW/weve-esi/staticdb/inner/sde"
)

func GetStation(id int32) (sde.Station, bool) {
	return sde.KVReaderStations.Get(id)
}
