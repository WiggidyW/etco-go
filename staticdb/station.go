package staticdb

import (
	"github.com/WiggidyW/eve-trading-co-go/staticdb/inner/sde"
)

func GetStation(id int32) (sde.Station, bool) {
	return sde.KVReaderStations.Get(id)
}
