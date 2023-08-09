package admin

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authing"
)

type AdminAccessType string

const (
	Read  AdminAccessType = "read"
	Write AdminAccessType = "write"
)

type AdminWriteParams struct {
	refreshToken string
	domain       string
	authList     AuthList
}

func (awcf AdminWriteParams) AuthRefreshToken() string {
	return awcf.refreshToken
}

type AuthingAdminWriteClient = authing.AuthingClient[
	AdminWriteParams,
	struct{},
	AdminWriteClient,
]

type AdminWriteClient struct {
	inner authing.AuthHashSetWriterClient
}

func (awc AdminWriteClient) Fetch(
	ctx context.Context,
	params AdminWriteParams,
) (*struct{}, error) {
	return awc.inner.Fetch(
		ctx,
		authing.AuthHashSetWriterParams{
			Key:         params.domain + "-" + string(Write),
			AuthHashSet: params.authList.toHashSet(),
		},
	)
}

type AdminReadParams struct {
	refreshToken string
	domain       string
}

func (arcf AdminReadParams) AuthRefreshToken() string {
	return arcf.refreshToken
}

type AuthingAdminReadClient = authing.AuthingClient[
	AdminReadParams,
	AuthList,
	AdminReadClient,
]

type AdminReadClient struct {
	inner authing.CachingAuthHashSetReaderClient
}

func (arc AdminReadClient) Fetch(
	ctx context.Context,
	params AdminReadParams,
) (*AuthList, error) {
	authHashSet, err := arc.inner.Fetch(
		ctx,
		authing.AuthHashSetReaderParams(
			params.domain+"-"+string(Read),
		),
	)
	if err != nil {
		return nil, err
	}
	authList := fromHashSet(authHashSet.Data())
	return &authList, nil
}

type AuthList struct {
	CharacterIDs   []int32
	CorporationIDs []int32
	AllianceIDs    []int32
}

func (al AuthList) toHashSet() authing.AuthHashSet {
	authHashSet := authing.AuthHashSet{
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

func fromHashSet(as authing.AuthHashSet) AuthList {
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
