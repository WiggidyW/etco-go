package staticdb

import (
	"github.com/WiggidyW/eve-trading-co-go/staticdb/inner/sde"
)

func GetSystem(id int32) (sde.System, bool) {
	return sde.KVReaderSystems.Get(id)
}
