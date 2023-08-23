package staticdb

import (
	"github.com/WiggidyW/weve-esi/staticdb/inner/sde"
)

func GetSystem(id int32) (sde.System, bool) {
	return sde.KVReaderSystems.Get(id)
}
