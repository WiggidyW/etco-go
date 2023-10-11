package primarysdedata

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

func (fetd FilterExtendedTypeData) addSDEGroup(
	sdeGroups FSDGroupIds,
	sdeCategories FSDCategoryIds,
	etcoGroups *[]b.Group,
	etcoGroupsIndexMap map[GroupId]int,
	etcoCategories *[]b.CategoryName,
	etcoCategoriesIndexMap map[CategoryId]int,
) (err error) {
	if index, err := getOrConvertSDEGroup(
		fetd.GroupId,
		sdeGroups,
		sdeCategories,
		etcoGroups,
		etcoGroupsIndexMap,
		etcoCategories,
		etcoCategoriesIndexMap,
	); err != nil {
		return err
	} else {
		// set the group index
		fetd.etcoTypeData.GroupIndex = index
		return nil
	}
}

func getOrConvertSDEGroup(
	groupId GroupId,
	sdeGroups FSDGroupIds,
	sdeCategories FSDCategoryIds,
	etcoGroups *[]b.Group,
	etcoGroupsIndexMap map[GroupId]int,
	etcoCategories *[]b.CategoryName,
	etcoCategoriesIndexMap map[CategoryId]int,
) (
	index int,
	err error,
) {
	// check if it was already added first
	index, exists := etcoGroupsIndexMap[groupId]
	if exists {
		return index, nil
	}

	// get the raw group, ensure that the group exists and is valid
	sdeGroup, exists := sdeGroups[groupId]
	if !exists {
		err = fmt.Errorf(
			"group id %d not found in sde",
			groupId,
		)
		return index, err
	} else if err = sdeGroup.validate(groupId); err != nil {
		return index, err
	}

	// find the category index
	categoryIndex, err := getOrConvertSDECategory(
		sdeGroup.CategoryId,
		sdeCategories,
		etcoCategories,
		etcoCategoriesIndexMap,
	)
	if err != nil {
		return index, err
	}

	// add to etcoGroups and etcoGroupsIndexMap
	index = len(*etcoGroups)
	etcoGroupsIndexMap[groupId] = index
	*etcoGroups = append(*etcoGroups, b.Group{
		Name:          sdeGroup.Name.En,
		CategoryIndex: categoryIndex,
	})

	return index, nil
}

func (gd GroupData) validate(id GroupId) error {
	if gd.CategoryId == 0 ||
		gd.Name.En == "" {
		return fmt.Errorf(
			"invalid group data: '%d': '%+v'",
			id,
			gd,
		)
	}
	return nil
}

func getOrConvertSDECategory(
	categoryId CategoryId,
	sdeCategories FSDCategoryIds,
	etcoCategories *[]b.CategoryName,
	etcoCategoriesIndexMap map[CategoryId]int,
) (
	index int,
	err error,
) {
	// check if it was already added first
	index, exists := etcoCategoriesIndexMap[categoryId]
	if exists {
		return index, nil
	}

	// get the raw category, ensure that the category exists and is valid
	sdeCategory, exists := sdeCategories[categoryId]
	if !exists {
		err = fmt.Errorf(
			"category id %d not found in sde",
			categoryId,
		)
		return index, err
	} else if err = sdeCategory.validate(categoryId); err != nil {
		return index, err
	}

	// add to etcoCategories and etcoCategoriesIndexMap
	index = len(*etcoCategories)
	etcoCategoriesIndexMap[categoryId] = index
	*etcoCategories = append(*etcoCategories, sdeCategory.Name.En)

	return index, nil
}

func (cd CategoryData) validate(id CategoryId) error {
	if cd.Name.En == "" {
		return fmt.Errorf(
			"invalid category data: '%d': '%+v'",
			id,
			cd,
		)
	}
	return nil
}
