package structureinfo

type StructureInfoModel struct {
	Name string `json:"name"`
	// OwnerId       int32             `json:"owner_id"`
	// Position      StructurePosition `json:"position"`
	SolarSystemId int32 `json:"solar_system_id"`
	// TypeId        *int32            `json:"type_id,omitempty"`
}

// type StructurePosition struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// 	Z float64 `json:"z"`
// }
