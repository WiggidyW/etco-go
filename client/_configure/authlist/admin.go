package admin

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

type AdminAccessType string

const (
	Read  AdminAccessType = "read"
	Write AdminAccessType = "write"
)

type AuthList struct {
	PermitCharacterIds   []int32
	BannedCharacterIds   []int32
	PermitCorporationIds []int32
	BannedCorporationIds []int32
	PermitAllianceIds    []int32
}

func (al AuthList) toHashSet() b.AuthHashSet {
	authHashSet := b.AuthHashSet{
		PermitCharacterIds: make(
			map[int32]struct{},
			len(al.PermitCharacterIds),
		),
		BannedCharacterIds: make(
			map[int32]struct{},
			len(al.BannedCharacterIds),
		),
		PermitCorporationIds: make(
			map[int32]struct{},
			len(al.PermitCorporationIds),
		),
		BannedCorporationIds: make(
			map[int32]struct{},
			len(al.BannedCorporationIds),
		),
		PermitAllianceIds: make(
			map[int32]struct{},
			len(al.PermitAllianceIds),
		),
	}
	for _, id := range al.PermitCharacterIds {
		authHashSet.PermitCharacterIds[id] = struct{}{}
	}
	for _, id := range al.BannedCharacterIds {
		authHashSet.BannedCharacterIds[id] = struct{}{}
	}
	for _, id := range al.PermitCorporationIds {
		authHashSet.PermitCorporationIds[id] = struct{}{}
	}
	for _, id := range al.BannedCorporationIds {
		authHashSet.BannedCorporationIds[id] = struct{}{}
	}
	for _, id := range al.PermitAllianceIds {
		authHashSet.PermitAllianceIds[id] = struct{}{}
	}
	return authHashSet
}

func fromHashSet(authHashSet b.AuthHashSet) (authList AuthList) {
	// benchmark shows that authList is faster than &authList for this
	authList = AuthList{
		PermitCharacterIds: make(
			[]int32,
			0,
			len(authHashSet.PermitCharacterIds),
		),
		BannedCharacterIds: make(
			[]int32,
			0,
			len(authHashSet.BannedCharacterIds),
		),
		PermitCorporationIds: make(
			[]int32,
			0,
			len(authHashSet.PermitCorporationIds),
		),
		BannedCorporationIds: make(
			[]int32,
			0,
			len(authHashSet.BannedCorporationIds),
		),
		PermitAllianceIds: make(
			[]int32,
			0,
			len(authHashSet.PermitAllianceIds),
		),
	}
	for id := range authHashSet.PermitCharacterIds {
		authList.PermitCharacterIds = append(
			authList.PermitCharacterIds,
			id,
		)
	}
	for id := range authHashSet.BannedCharacterIds {
		authList.BannedCharacterIds = append(
			authList.BannedCharacterIds,
			id,
		)
	}
	for id := range authHashSet.PermitCorporationIds {
		authList.PermitCorporationIds = append(
			authList.PermitCorporationIds,
			id,
		)
	}
	for id := range authHashSet.BannedCorporationIds {
		authList.BannedCorporationIds = append(
			authList.BannedCorporationIds,
			id,
		)
	}
	for id := range authHashSet.PermitAllianceIds {
		authList.PermitAllianceIds = append(
			authList.PermitAllianceIds,
			id,
		)
	}
	return authList
}
