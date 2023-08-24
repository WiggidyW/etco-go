package contracts

import (
	"time"

	"github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	cc "github.com/WiggidyW/eve-trading-co-go/client/esi/model/contractscorporation"
)

type Contracts struct {
	ShopContracts    map[string]Contract
	BuybackContracts map[string]Contract
}

func newContracts() *Contracts {
	return &Contracts{
		ShopContracts:    make(map[string]Contract),
		BuybackContracts: make(map[string]Contract),
	}
}

func (c *Contracts) filterAddEntry(
	corporationId int32,
	entry cc.ContractsCorporationEntry,
) {
	// filter out contracts that aren't buyback or shop contracts
	if (entry.AssigneeId != corporationId &&
		entry.IssuerCorporationId != corporationId) ||
		entry.Title == nil ||
		*entry.Title == "" ||
		entry.Type != "item_exchange" ||
		entry.EndLocationId == nil ||
		entry.Price == nil {
		return
	}

	// ensure that the code type matches the contract direction
	code, codeType := appraisal.ParseCode(*entry.Title)

	// buyback requires user -> corp
	if codeType == appraisal.BuybackCode { // BuybackCode
		if entry.AssigneeId == corporationId {
			c.BuybackContracts[code] = cFromEntry(entry)
		}

		// shop requires corp -> user
	} else if codeType == appraisal.ShopCode { // ShopCode
		if entry.IssuerCorporationId == corporationId {
			c.ShopContracts[code] = cFromEntry(entry)
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

func cFromEntry(entry cc.ContractsCorporationEntry) Contract {
	return Contract{
		ContractId:   entry.ContractId,
		Status:       sFromString(entry.Status),
		Issued:       entry.DateIssued,
		Expires:      entry.DateExpired,
		LocationId:   *entry.EndLocationId,
		Price:        *entry.Price,
		HasReward:    entry.Reward != nil && *entry.Reward > 1,
		IssuerCorpId: entry.IssuerCorporationId,
		IssuerCharId: entry.IssuerId,
		AssigneeId:   entry.AssigneeId,
		AssigneeType: atFromString(entry.Availability),
	}
}
