package contractscorporation

import "time"

type ContractsCorporationEntry struct {
	// AcceptorId   int32  `json:"acceptor_id"`
	AssigneeId   int32  `json:"assignee_id"`
	Availability string `json:"availability"`
	// Buyout              *float64   `json:"buyout,omitempty"`
	// Collateral          *float64   `json:"collateral,omitempty"`
	ContractId int32 `json:"contract_id"`
	// DateAccepted        *time.Time `json:"date_accepted,omitempty"`
	// DateCompleted       *time.Time `json:"date_completed,omitempty"`
	DateExpired time.Time `json:"date_expired"`
	DateIssued  time.Time `json:"date_issued"`
	// DaysToComplete      *int32     `json:"days_to_complete,omitempty"`
	EndLocationId *int64 `json:"end_location_id,omitempty"`
	// ForCorporation      bool       `json:"for_corporation"`
	IssuerCorporationId int32    `json:"issuer_corporation_id"`
	IssuerId            int32    `json:"issuer_id"`
	Price               *float64 `json:"price,omitempty"`
	Reward              *float64 `json:"reward,omitempty"`
	// StartLocationId     *int64     `json:"start_location_id,omitempty"`
	Status string   `json:"status"`
	Title  *string  `json:"title,omitempty"`
	Type   string   `json:"type"`
	Volume *float64 `json:"volume,omitempty"`
}
