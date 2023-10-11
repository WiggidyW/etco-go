package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBPurchaseQueueParams struct {
	TypeNamingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap]
	QueueInclude      PurchaseQueueInclude
}

type PBContractQueueParams struct {
	TypeNamingSession   *staticdb.TypeNamingSession[*staticdb.SyncIndexMap]
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
	QueueInclude        ContractQueueInclude
}

type PBStatusAppraisalParams struct {
	TypeNamingSession   *staticdb.TypeNamingSession[*staticdb.LocalIndexMap]
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker]
	AppraisalCode       string
	StatusInclude       AppraisalStatusInclude
}

type ContractQueueInclude uint8

const (
	// items, code appraisal, new appraisal

	//
	CQI_NONE ContractQueueInclude = iota

	//   T
	CQI_ITEMS

	//                 T
	CQI_CODE_APPRAISAL

	//   T             T
	CQI_ITEMS_AND_CODE_APPRAISAL

	//                 T               T
	CQI_CODE_APPRAISAL_AND_NEW_APPRAISAL

	//   T             T               T
	CQI_ITEMS_AND_CODE_APPRAISAL_AND_NEW_APPRAISAL
)

func NewContractQueueInclude(
	include_items bool,
	include_code_appraisal bool,
	include_new_appraisal bool,
) ContractQueueInclude {
	if !include_items && !include_code_appraisal && !include_new_appraisal { // 000
		return CQI_NONE
	} else if !include_items && !include_code_appraisal && include_new_appraisal { // 001
		return CQI_NONE
	} else if !include_items && include_code_appraisal && !include_new_appraisal { // 010
		return CQI_CODE_APPRAISAL
	} else if !include_items && include_code_appraisal && include_new_appraisal { // 011
		return CQI_CODE_APPRAISAL
	} else if include_items && !include_code_appraisal && !include_new_appraisal { // 100
		return CQI_ITEMS
	} else if include_items && !include_code_appraisal && include_new_appraisal { // 101
		return CQI_ITEMS
	} else if include_items && include_code_appraisal && !include_new_appraisal { // 110
		return CQI_ITEMS_AND_CODE_APPRAISAL
	} else if include_items && include_code_appraisal && include_new_appraisal { // 111
		return CQI_ITEMS_AND_CODE_APPRAISAL_AND_NEW_APPRAISAL
	}
	panic("unreachable")
}

type PurchaseQueueInclude uint8

const (
	// code appraisal, new appraisal

	//
	PQI_NONE PurchaseQueueInclude = iota

	//       T
	PQI_CODE_APPRAISAL

	//       T               T
	PQI_CODE_APPRAISAL_AND_NEW_APPRAISAL
)

func NewPurchaseQueueInclude(
	include_code_appraisal bool,
	include_new_appraisal bool,
) PurchaseQueueInclude {
	if !include_code_appraisal && !include_new_appraisal { // 00
		return PQI_NONE
	} else if !include_code_appraisal && include_new_appraisal { // 01
		return PQI_NONE
	} else if include_code_appraisal && !include_new_appraisal { // 10
		return PQI_CODE_APPRAISAL
	} else if include_code_appraisal && include_new_appraisal { // 11
		return PQI_CODE_APPRAISAL_AND_NEW_APPRAISAL
	}
	panic("unreachable")
}

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
	if !include_items && !include_location_info { // 00
		return ASI_NONE
	} else if !include_items && include_location_info { // 01
		return ASI_LOCATION_INFO
	} else if include_items && !include_location_info { // 10
		return ASI_ITEMS
	} else if include_items && include_location_info { // 11
		return ASI_ITEMS_AND_LOCATION_INFO
	}
	panic("unreachable")
}
