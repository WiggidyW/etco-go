package contracts

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
)

type Contracts struct {
	ShopContracts    map[string]Contract
	BuybackContracts map[string]Contract
	HaulContracts    map[string]Contract
}

func newContracts() Contracts {
	return Contracts{
		ShopContracts:    make(map[string]Contract),
		BuybackContracts: make(map[string]Contract),
		HaulContracts:    make(map[string]Contract),
	}
}

func (c *Contracts) filterAddEntries(
	entries []esi.ContractsEntry,
) {
	for _, entry := range entries {
		c.filterAddEntry(entry)
	}
}

func (c *Contracts) filterAddEntry(
	entry esi.ContractsEntry,
) {
	// filter out contracts that aren't buyback or shop contracts
	if (entry.AssigneeId != build.CORPORATION_ID &&
		entry.IssuerCorporationId != build.CORPORATION_ID) ||
		entry.Title == nil || *entry.Title == "" ||
		(entry.Type != "item_exchange" && entry.Type != "courier") ||
		entry.EndLocationId == nil ||
		entry.Price == nil {
		return
	}

	// ensure that the code type matches the contract direction
	code, codeType := appraisalcode.ParseCode(*entry.Title)

	var contracts map[string]Contract
	var entitySource ContractEntitySource
	switch codeType {
	case appraisalcode.BuybackCode:
		if entry.Type == "item_exchange" &&
			entry.AssigneeId == build.CORPORATION_ID {
			contracts = c.BuybackContracts
			entitySource = CESIssuer
		} else {
			return
		}
	case appraisalcode.ShopCode:
		if entry.Type == "item_exchange" &&
			entry.IssuerCorporationId == build.CORPORATION_ID {
			contracts = c.ShopContracts
			entitySource = CESAssignee
		} else {
			return
		}
	case appraisalcode.HaulCode:
		if entry.Type == "courier" &&
			entry.AssigneeId == build.CORPORATION_ID {
			contracts = c.HaulContracts
			entitySource = CESIssuer
		} else {
			return
		}
	default: // UnknownCode (or invalid enum value)
		return
	}

	existing, ok := contracts[code]
	if !ok || entry.DateIssued.After(existing.Issued) {
		contracts[code] = fromEntry(entry, entitySource)
	}
}

type Contract struct {
	ContractId      int32
	Status          Status
	Issued          time.Time
	Expires         time.Time
	StartLocationId int64 // 0 unless Haul
	LocationId      int64
	Price           float64 // collateral if haul
	Reward          float64
	EntityKind      EntityKind
	EntityId        int32
}

type ContractEntitySource bool

const (
	CESIssuer   ContractEntitySource = false
	CESAssignee ContractEntitySource = true
)

func fromEntry(
	entry esi.ContractsEntry,
	entitySource ContractEntitySource,
) Contract {
	var startLocationId int64
	if entry.StartLocationId != nil {
		startLocationId = *entry.StartLocationId
	}
	var price float64
	if entry.Collateral != nil && *entry.Collateral > 0.0 {
		price = *entry.Collateral
	} else if entry.Price != nil {
		price = *entry.Price
	}
	var reward float64
	if entry.Reward != nil {
		reward = *entry.Reward
	}
	var entityKind EntityKind
	var entityId int32
	if entitySource == CESIssuer {
		entityKind = EKCharacter
		entityId = entry.IssuerId
	} else /* if entitySource == CESAssignee */ {
		entityKind = entityKindFromStr(entry.Availability)
		entityId = entry.AssigneeId
	}
	return Contract{
		ContractId:      entry.ContractId,
		Status:          sFromString(entry.Status),
		Issued:          entry.DateIssued,
		Expires:         entry.DateExpired,
		StartLocationId: startLocationId,
		LocationId:      *entry.EndLocationId,
		Price:           price,
		Reward:          reward,
		EntityKind:      entityKind,
		EntityId:        entityId,
	}
}

func (c Contract) ToProto(
	startLocationInfo *proto.LocationInfo,
	locationInfo *proto.LocationInfo,
) *proto.Contract {
	return &proto.Contract{
		ContractId:        c.ContractId,
		Status:            c.Status.ToProto(),
		Issued:            c.Issued.Unix(),
		Expires:           c.Expires.Unix(),
		StartLocationInfo: startLocationInfo,
		LocationInfo:      locationInfo,
		Price:             c.Price,
		Reward:            c.Reward,
		EntityKind:        c.EntityKind.ToProto(),
		EntityId:          c.EntityId,
	}
}
