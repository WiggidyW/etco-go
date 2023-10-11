package primarysdedata

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

type FilterExtendedTypeData struct {
	etcoTypeData  *b.TypeData
	TypeId        b.TypeId
	GroupId       GroupId
	Volume        float64
	MarketGroupId MarketGroupId
	Name          TypeNames
	PortionSize   *float64
}

// keep only published types with a non-nil market group id
func filterExtendSDETypeDataMap(
	sdeTypeDataMap FSDTypeIds,
) (
	filterExtendedTypeDatas []FilterExtendedTypeData,
	err error,
) {
	filterExtendedTypeDatas = make([]FilterExtendedTypeData, 0)

	for typeId, sdeTypeData := range sdeTypeDataMap {
		if !sdeTypeData.Published || sdeTypeData.MarketGroupId == nil {
			continue
		}
		if err := sdeTypeData.validate(); err != nil {
			return nil, err
		}
		filterExtendedTypeData := FilterExtendedTypeData{
			etcoTypeData: &b.TypeData{
				Name: sdeTypeData.Name.En,
				// ReprocessedMaterials: []b.ReprocessedMaterial{},
				GroupIndex:       -1,
				MarketGroupIndex: -1,
				VolumeIndex:      -1,
			},
			TypeId:        typeId,
			GroupId:       *sdeTypeData.GroupId,
			Volume:        sdeTypeData.Volume,
			MarketGroupId: *sdeTypeData.MarketGroupId,
			Name:          sdeTypeData.Name,
			PortionSize:   sdeTypeData.PortionSize,
		}
		filterExtendedTypeDatas = append(
			filterExtendedTypeDatas,
			filterExtendedTypeData,
		)
	}

	return filterExtendedTypeDatas, nil
}

func (td TypeData) validate() error {
	if td.GroupId == nil || td.Name.En == "" {
		return fmt.Errorf(
			"invalid type data: '%+v'",
			td,
		)
	}
	return nil
}
