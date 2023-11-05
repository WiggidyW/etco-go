package bucket

import (
	"context"
	"time"

	b "github.com/WiggidyW/etco-go-bucket"
)

func GetAuthList(
	ctx context.Context,
	domain string,
) (
	rep AuthList,
	expires *time.Time,
	err error,
) {
	var authHashSet b.AuthHashSet
	authHashSet, expires, err = GetAuthHashSet(ctx, domain)
	if err == nil {
		rep = authListfromAHS(authHashSet)
	}
	return rep, expires, err
}

func SetAuthList(
	ctx context.Context,
	domain string,
	rep AuthList,
) (
	err error,
) {
	return SetAuthHashSet(ctx, domain, authListToAHS(rep))
}

type AuthList struct {
	PermitCharacterIds   []int32
	BannedCharacterIds   []int32
	PermitCorporationIds []int32
	BannedCorporationIds []int32
	PermitAllianceIds    []int32
}

func authListToAHS(authList AuthList) (authHashSet b.AuthHashSet) {
	authHashSet = b.AuthHashSet{
		PermitCharacterIds: make(
			map[int32]struct{},
			len(authList.PermitCharacterIds),
		),
		BannedCharacterIds: make(
			map[int32]struct{},
			len(authList.BannedCharacterIds),
		),
		PermitCorporationIds: make(
			map[int32]struct{},
			len(authList.PermitCorporationIds),
		),
		BannedCorporationIds: make(
			map[int32]struct{},
			len(authList.BannedCorporationIds),
		),
		PermitAllianceIds: make(
			map[int32]struct{},
			len(authList.PermitAllianceIds),
		),
	}
	for _, id := range authList.PermitCharacterIds {
		authHashSet.PermitCharacterIds[id] = struct{}{}
	}
	for _, id := range authList.BannedCharacterIds {
		authHashSet.BannedCharacterIds[id] = struct{}{}
	}
	for _, id := range authList.PermitCorporationIds {
		authHashSet.PermitCorporationIds[id] = struct{}{}
	}
	for _, id := range authList.BannedCorporationIds {
		authHashSet.BannedCorporationIds[id] = struct{}{}
	}
	for _, id := range authList.PermitAllianceIds {
		authHashSet.PermitAllianceIds[id] = struct{}{}
	}
	return authHashSet
}

func authListfromAHS(authHashSet b.AuthHashSet) (authList AuthList) {
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
