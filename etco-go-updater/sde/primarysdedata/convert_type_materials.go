package primarysdedata

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

// inserts the reprocessed materials into extendedTypeDatas
func (fetd FilterExtendedTypeData) addSDETypeMaterials(
	sdeTypeMaterialsMap FSDTypeMaterials,
) (err error) {
	fetd.etcoTypeData.ReprocessedMaterials, err = convertSDETypeMaterials(
		sdeTypeMaterialsMap[fetd.TypeId].Materials,
		fetd.PortionSize,
	)
	return err
}

func (tm TypeMaterial) validate() error {
	if tm.MaterialTypeId == 0 ||
		tm.Quantity == 0 {
		return fmt.Errorf(
			"invalid type material data: '%+v'",
			tm,
		)
	}
	return nil
}

func convertSDETypeMaterials(
	sdeTypeMaterials []TypeMaterial,
	portionSizePtr *float64,
) (
	etcoTypeMaterials []b.ReprocessedMaterial,
	err error,
) {
	if len(sdeTypeMaterials) == 0 {
		return []b.ReprocessedMaterial{}, nil
	}

	var portionSize float64
	if portionSizePtr != nil {
		portionSize = *portionSizePtr
	} else {
		portionSize = 1.0
	}

	etcoTypeMaterials = make(
		[]b.ReprocessedMaterial,
		0,
		len(sdeTypeMaterials),
	)
	for _, sdeTypeMaterial := range sdeTypeMaterials {
		if err = sdeTypeMaterial.validate(); err != nil {
			return nil, err
		}
		etcoTypeMaterials = append(
			etcoTypeMaterials,
			b.ReprocessedMaterial{
				TypeId:   sdeTypeMaterial.MaterialTypeId,
				Quantity: sdeTypeMaterial.Quantity / portionSize,
			},
		)
	}

	return etcoTypeMaterials, nil
}
