package primarysdedata

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

// inserts the reprocessed materials into extendedTypeDatas
func (fetd FilterExtendedTypeData) addSDETypeVolume(
	etcoTypeVolumes *[]b.TypeVolume,
	etcoTypeVolumesIndexMap map[float64]int,
) {
	index, exists := etcoTypeVolumesIndexMap[fetd.Volume]
	if !exists {
		index = len(*etcoTypeVolumes)
		etcoTypeVolumesIndexMap[fetd.Volume] = index
		*etcoTypeVolumes = append(
			*etcoTypeVolumes,
			fetd.Volume,
		)
	}
	fetd.etcoTypeData.VolumeIndex = index
}
