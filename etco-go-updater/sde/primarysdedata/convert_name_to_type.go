package primarysdedata

import (
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

// First type ID will be used
var HANDLE_NAME_CONFLICTS = map[string][]b.TypeId{
	"Enforcer - Frigate Crate":                         {63780, 63782},
	"Stabber Glacial Drift SKIN":                       {44171, 46894},
	"Scythe Glacial Drift SKIN":                        {44169, 46893},
	"Badger Wiyrkomi SKIN":                             {60106, 36333},
	"Tengu Ultra Jungle":                               {48543, 48544},
	"Brutix Serpentis SKIN":                            {39584, 42177},
	"Festival Facial Augmentation and Snowballs Crate": {53493, 53513},
	"Catalyst Serpentis SKIN":                          {39585, 42162},
	"Intaki Emerald Metallic - Limited":                {84014, 84088, 84089, 84090, 84091, 84092, 84093, 84094, 84095, 84096, 84097},
	"Spiked Quafe":                                     {21661, 54165},
	"Festival Calm Dark Filament and Snowballs Crate":  {53489, 53492},
	"Festival Skill Points and Snowballs Crate":        {53497, 53499, 53507, 53510, 53515, 53517},
	"Minmatar Liberation Day Apparel Crate":            {63780, 63782},
}

var NTTI_ORDER = []func(tn TypeNames) string{
	// primary will ensure that no 2 types have the same name.
	// If a name is already present, it will be skipped unless primary
	// For primary, duplicate names will raise an error
	func(tn TypeNames) string { return tn.En }, // 1st / primary
	func(tn TypeNames) string { return tn.Ru }, // 2nd
	func(tn TypeNames) string { return tn.De }, // 3rd
	func(tn TypeNames) string { return tn.Fr }, // 4th
	func(tn TypeNames) string { return tn.Zh }, // 5th
	func(tn TypeNames) string { return tn.Ja }, // 6th
	func(tn TypeNames) string { return tn.Es }, // 7th
	func(tn TypeNames) string { return tn.It }, // 8th
}

func transceiveConvertFETypeDatasToETCONameToTypeId(
	feTypeDatas []FilterExtendedTypeData,
	chnSend chanresult.ChanSendResult[map[b.TypeName]b.TypeId],
) error {
	if etcoNameToTypeId, err := convertFETypeDatasToETCONameToTypeId(
		feTypeDatas,
	); err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(etcoNameToTypeId)
	}
}

func convertFETypeDatasToETCONameToTypeId(
	feTypeDatas []FilterExtendedTypeData,
) (
	etcoNameToTypeId map[b.TypeName]b.TypeId,
	err error,
) {
	etcoNameToTypeId = make(map[b.TypeName]b.TypeId)

	for _, feTypeData := range feTypeDatas {
		name := NTTI_ORDER[0](feTypeData.Name)
		if err := addNameToNameToTypeId(
			etcoNameToTypeId,
			feTypeData.TypeId,
			name,
			false,
		); err != nil {
			return nil, err
		}
	}

	for _, nameFunc := range NTTI_ORDER[1:] {
		for _, feTypeData := range feTypeDatas {
			name := nameFunc(feTypeData.Name)
			_ = addNameToNameToTypeId(
				etcoNameToTypeId,
				feTypeData.TypeId,
				name,
				true,
			)
		}
	}

	return etcoNameToTypeId, nil
}

func addNameToNameToTypeId(
	nameToTypeId map[b.TypeName]b.TypeId,
	newTypeId b.TypeId,
	name string,
	conflictOkay bool,
) error {
	typeId, present := nameToTypeId[name]
	if !present {
		nameToTypeId[name] = newTypeId
	} else {
		handleConflict, ok := HANDLE_NAME_CONFLICTS[name]
		if ok && contains(handleConflict, typeId) {
			nameToTypeId[name] = handleConflict[0]
		} else if !conflictOkay {
			return fmt.Errorf(
				"name conflict for '%s' - '%d' and '%d'",
				name,
				typeId,
				newTypeId,
			)
		}
	}
	return nil
}

func contains(arr []b.TypeId, val b.TypeId) bool {
	for _, arrVal := range arr {
		if arrVal == val {
			return true
		}
	}
	return false
}
