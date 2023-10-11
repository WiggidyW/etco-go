package etcogobucket

type AuthHashSet struct {
	// checked in the following order, breaking on first match
	BannedCharacterIds   map[int32]struct{} // if banned, access not given
	PermitCharacterIds   map[int32]struct{} // if permitted, access given
	BannedCorporationIds map[int32]struct{} // if banned, access not given
	PermitCorporationIds map[int32]struct{} // if permitted, access given
	// BannedAllianceIds    map[int32]struct{} // if banned, access not given
	PermitAllianceIds map[int32]struct{} // if permitted, access given
}

func (ahs AuthHashSet) BannedCharacter(id int32) bool {
	_, ok := ahs.BannedCharacterIds[id]
	return ok
}

func (ahs AuthHashSet) PermittedCharacter(id int32) bool {
	_, ok := ahs.PermitCharacterIds[id]
	return ok
}

func (ahs AuthHashSet) BannedCorporation(id int32) bool {
	_, ok := ahs.BannedCorporationIds[id]
	return ok
}

func (ahs AuthHashSet) PermittedCorporation(id int32) bool {
	_, ok := ahs.PermitCorporationIds[id]
	return ok
}

// func (ahs AuthHashSet) BannedAlliance(id int32) bool {
// 	_, ok := ahs.BannedAllianceIds[id]
// 	return ok
// }

func (ahs AuthHashSet) PermittedAlliance(id int32) bool {
	_, ok := ahs.PermitAllianceIds[id]
	return ok
}
