package primarysdedata

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

// converts sdeMarketGroups to etcoMarketGroups and inserts their indexes into extendedTypeDatas
func (fetd FilterExtendedTypeData) addSDEMarketGroup(
	sdeMarketGroups FSDMarketGroups,
	etcoMarketGroups *[]b.MarketGroup,
	etcoMarketGroupsIndexMap map[MarketGroupId]int,
) (err error) {
	// convert and add the market group
	if _, index, err := getOrConvertSDEMarketGroup(
		fetd.MarketGroupId,
		sdeMarketGroups,
		etcoMarketGroups,
		etcoMarketGroupsIndexMap,
	); err != nil {
		return err
	} else {
		// set the market group index
		fetd.etcoTypeData.MarketGroupIndex = index
		return nil
	}

}

func getOrConvertSDEMarketGroup(
	marketGroupId MarketGroupId,
	sdeMarketGroups FSDMarketGroups,
	etcoMarketGroups *[]b.MarketGroup,
	etcoMarketGroupsIndexMap map[MarketGroupId]int,
) (
	numParents uint8,
	index int,
	err error,
) {
	// check if it was already added first
	index, exists := etcoMarketGroupsIndexMap[marketGroupId]
	if exists {
		return (*etcoMarketGroups)[index].NumParents, index, nil
	}

	// ensure that the market group exists and is valid
	sdeMarketGroup, exists := sdeMarketGroups[marketGroupId]
	if !exists {
		err = fmt.Errorf(
			"market group id %d not found in sde",
			marketGroupId,
		)
		return numParents, index, err
	} else if err = sdeMarketGroup.validate(marketGroupId); err != nil {
		return numParents, index, err
	}

	// find num parents and parent index
	var parentIndex int
	if sdeMarketGroup.ParentGroupId != nil {
		// if we have a parent, we need to find its number of parents
		numParents, parentIndex, err = getOrConvertSDEMarketGroup(
			*sdeMarketGroup.ParentGroupId,
			sdeMarketGroups,
			etcoMarketGroups,
			etcoMarketGroupsIndexMap,
		)
		if err != nil {
			return numParents, index, err
		}
		// our number of parents is our parent's number of parents + 1
		numParents = numParents + 1
	} else {
		// if we have no parents, our number of parents is 0 and our parent index is -1 (no parent)
		numParents = 0
		parentIndex = -1
	}

	// add the market group to the list and update the index map
	index = len(*etcoMarketGroups)
	etcoMarketGroupsIndexMap[marketGroupId] = index
	*etcoMarketGroups = append(*etcoMarketGroups, b.MarketGroup{
		Name:        sdeMarketGroup.NameId.En,
		ParentIndex: parentIndex,
		NumParents:  numParents,
	})

	return numParents, index, nil
}

func (mgd MarketGroupData) validate(id MarketGroupId) error {
	if mgd.NameId.En == "" {
		return fmt.Errorf(
			"invalid market group data: '%d': '%+v'",
			id,
			mgd,
		)
	}
	return nil
}
