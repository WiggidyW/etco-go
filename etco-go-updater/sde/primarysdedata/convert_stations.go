package primarysdedata

import (
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

func transceiveConvertSDEStations(
	sdeStations BSDStaStations,
	chnSend chanresult.ChanSendResult[map[b.StationId]b.Station],
) error {
	if etcoStations, err := convertSDEStations(sdeStations); err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(etcoStations)
	}
}

func convertSDEStations(
	sdeStations BSDStaStations,
) (
	etcoStations map[b.StationId]b.Station,
	err error,
) {
	etcoStations = make(map[b.StationId]b.Station, len(sdeStations))

	for _, sdeStation := range sdeStations {
		if err := sdeStation.validate(); err != nil {
			return nil, err
		}
		etcoStations[sdeStation.StationId] = b.Station{
			Name:     sdeStation.Name,
			SystemId: sdeStation.SystemId,
		}
	}

	return etcoStations, nil
}

func (ss StaStation) validate() error {
	if ss.StationId == 0 ||
		ss.Name == "" ||
		ss.SystemId == 0 {
		return fmt.Errorf("invalid station data: %+v", ss)
	}
	return nil
}
