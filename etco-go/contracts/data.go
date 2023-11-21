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
}

func newContracts() Contracts {
	return Contracts{
		ShopContracts:    make(map[string]Contract),
		BuybackContracts: make(map[string]Contract),
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
		entry.Title == nil ||
		*entry.Title == "" ||
		entry.Type != "item_exchange" ||
		entry.EndLocationId == nil ||
		entry.Price == nil {
		return
	}

	// ensure that the code type matches the contract direction
	code, codeType := appraisalcode.ParseCode(*entry.Title)

	// buyback requires user -> corp
	if codeType == appraisalcode.BuybackCode { // BuybackCode
		if entry.AssigneeId == build.CORPORATION_ID {
			existing, ok := c.BuybackContracts[code]
			if !ok || entry.DateIssued.After(existing.Issued) {
				c.BuybackContracts[code] = fromEntry(entry)
			}
		}

		// shop requires corp -> user
	} else if codeType == appraisalcode.ShopCode { // ShopCode
		if entry.IssuerCorporationId == build.CORPORATION_ID {
			existing, ok := c.ShopContracts[code]
			if !ok || entry.DateIssued.After(existing.Issued) {
				c.ShopContracts[code] = fromEntry(entry)
			}
		}

	} // else UnknownCode
}

type Contract struct {
	ContractId   int32
	Status       Status
	Issued       time.Time
	Expires      time.Time
	LocationId   int64
	Price        float64
	HasReward    bool
	IssuerCorpId int32
	IssuerCharId int32
	AssigneeId   int32
	AssigneeType AssigneeType
}

func fromEntry(entry esi.ContractsEntry) Contract {
	price := *entry.Price
	if entry.Reward != nil {
		price -= *entry.Reward
	}
	return Contract{
		ContractId:   entry.ContractId,
		Status:       sFromString(entry.Status),
		Issued:       entry.DateIssued,
		Expires:      entry.DateExpired,
		LocationId:   *entry.EndLocationId,
		Price:        price,
		HasReward:    entry.Reward != nil && *entry.Reward > 1,
		IssuerCorpId: entry.IssuerCorporationId,
		IssuerCharId: entry.IssuerId,
		AssigneeId:   entry.AssigneeId,
		AssigneeType: atFromString(entry.Availability),
	}
}

func (c Contract) ToProto(
	locationInfo *proto.LocationInfo,
) *proto.Contract {
	return &proto.Contract{
		ContractId:   c.ContractId,
		Status:       c.Status.ToProto(),
		Issued:       c.Issued.Unix(),
		Expires:      c.Expires.Unix(),
		LocationInfo: locationInfo,
		Price:        c.Price,
		IssuerCorpId: c.IssuerCorpId,
		IssuerCharId: c.IssuerCharId,
		AssigneeId:   c.AssigneeId,
		AssigneeType: c.AssigneeType.ToProto(),
	}
}
