package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

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
	if !include_items && !include_code_appraisal && !include_new_appraisal {
		return CQI_NONE
	} else if !include_items && !include_code_appraisal && include_new_appraisal {
		return CQI_NONE
	} else if !include_items && include_code_appraisal && !include_new_appraisal {
		return CQI_CODE_APPRAISAL
	} else if !include_items && include_code_appraisal && include_new_appraisal {
		return CQI_CODE_APPRAISAL
	} else if include_items && !include_code_appraisal && !include_new_appraisal {
		return CQI_ITEMS
	} else if include_items && !include_code_appraisal && include_new_appraisal {
		return CQI_ITEMS
	} else if include_items && include_code_appraisal && !include_new_appraisal {
		return CQI_ITEMS_AND_CODE_APPRAISAL
	} else if include_items && include_code_appraisal && include_new_appraisal {
		return CQI_ITEMS_AND_CODE_APPRAISAL_AND_NEW_APPRAISAL
	}
	panic("unreachable")
}

type PBContractQueueParams struct {
	TypeNamingSession   *staticdb.TypeNamingSession[*staticdb.SyncIndexMap]
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
	QueueInclude        ContractQueueInclude
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
	if !include_code_appraisal && !include_new_appraisal {
		return PQI_NONE
	} else if !include_code_appraisal && include_new_appraisal {
		return PQI_NONE
	} else if include_code_appraisal && !include_new_appraisal {
		return PQI_CODE_APPRAISAL
	} else if include_code_appraisal && include_new_appraisal {
		return PQI_CODE_APPRAISAL_AND_NEW_APPRAISAL
	}
	panic("unreachable")
}

type PBPurchaseQueueParams struct {
	TypeNamingSession   *staticdb.TypeNamingSession[*staticdb.SyncIndexMap]
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
	QueueInclude        PurchaseQueueInclude
}
