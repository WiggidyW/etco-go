package admin

import (
	"github.com/WiggidyW/weve-esi/client/configure/authhashset"
)

type AdminAccessType string

const (
	Read  AdminAccessType = "read"
	Write AdminAccessType = "write"
)

type AuthList struct {
	CharacterIDs   []int32
	CorporationIDs []int32
	AllianceIDs    []int32
}

func (al AuthList) toHashSet() authhashset.AuthHashSet {
	authHashSet := authhashset.AuthHashSet{
		CharacterIDs: make(
			map[int32]struct{},
			len(al.CharacterIDs),
		),
		CorporationIDs: make(
			map[int32]struct{},
			len(al.CorporationIDs),
		),
		AllianceIDs: make(
			map[int32]struct{},
			len(al.AllianceIDs),
		),
	}
	for _, id := range al.CharacterIDs {
		authHashSet.CharacterIDs[id] = struct{}{}
	}
	for _, id := range al.CorporationIDs {
		authHashSet.CorporationIDs[id] = struct{}{}
	}
	for _, id := range al.AllianceIDs {
		authHashSet.AllianceIDs[id] = struct{}{}
	}
	return authHashSet
}

func fromHashSet(as authhashset.AuthHashSet) AuthList {
	// benchmark shows that authList is faster than &authList for this
	authList := AuthList{
		CharacterIDs:   make([]int32, 0, len(as.CharacterIDs)),
		CorporationIDs: make([]int32, 0, len(as.CorporationIDs)),
		AllianceIDs:    make([]int32, 0, len(as.AllianceIDs)),
	}
	for id := range as.CharacterIDs {
		authList.CharacterIDs = append(authList.CharacterIDs, id)
	}
	for id := range as.CorporationIDs {
		authList.CorporationIDs = append(authList.CorporationIDs, id)
	}
	for id := range as.AllianceIDs {
		authList.AllianceIDs = append(authList.AllianceIDs, id)
	}
	return authList
}
