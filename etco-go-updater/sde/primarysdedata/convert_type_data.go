package primarysdedata

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

func convertFETypeDatasToETCOTypeDataMap(
	feTypeDatas []FilterExtendedTypeData,
) (
	etcoTypeDataMap map[b.TypeId]b.TypeData,
	err error,
) {
	etcoTypeDataMap = make(map[b.TypeId]b.TypeData, len(feTypeDatas))

	for _, feTypeData := range feTypeDatas {
		if err := feTypeData.validate(); err != nil {
			return nil, err
		}
		etcoTypeDataMap[feTypeData.TypeId] = *feTypeData.etcoTypeData
	}

	return etcoTypeDataMap, nil
}

func (fetd FilterExtendedTypeData) validate() error {
	if fetd.TypeId == 0 ||
		fetd.etcoTypeData == nil ||
		fetd.etcoTypeData.Name == "" ||
		fetd.etcoTypeData.GroupIndex == -1 ||
		fetd.etcoTypeData.MarketGroupIndex == -1 ||
		fetd.etcoTypeData.VolumeIndex == -1 {
		return fmt.Errorf(
			"invalid filter extended type data: '%+v'",
			fetd,
		)
	}
	return nil
}
