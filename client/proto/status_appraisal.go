package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type AppraisalStatusInclude uint8

const (
	// contract items, location info

	//
	ASI_NONE AppraisalStatusInclude = iota

	//       T
	ASI_ITEMS

	//                       T
	ASI_LOCATION_INFO

	//       T               T
	ASI_ITEMS_AND_LOCATION_INFO
)

func NewAppraisalStatusInclude(
	include_items bool,
	include_location_info bool,
) AppraisalStatusInclude {
	if !include_items && !include_location_info {
		return ASI_NONE
	} else if !include_items && include_location_info {
		return ASI_LOCATION_INFO
	} else if include_items && !include_location_info {
		return ASI_ITEMS
	} else if include_items && include_location_info {
		return ASI_ITEMS_AND_LOCATION_INFO
	}
	panic("unreachable")
}

type PBStatusAppraisalParams struct {
	TypeNamingSession   *staticdb.TypeNamingSession[*staticdb.LocalIndexMap]
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker]
	AppraisalCode       string
	StatusInclude       AppraisalStatusInclude
}
