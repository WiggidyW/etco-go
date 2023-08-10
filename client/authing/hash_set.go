package authing

import (
	"github.com/WiggidyW/weve-esi/client/storage"
)

type AuthHashSetReaderParams = storage.BucketReaderParams

type AuthHashSetReaderClient = storage.CachingBucketReaderClient[AuthHashSet]

type AuthHashSetWriterParams = storage.BucketWriterParams[AuthHashSet]

type AuthHashSetWriterClient = storage.AntiCachingBucketWriterClient[AuthHashSet]

type AuthHashSet struct {
	CharacterIDs   map[int32]struct{}
	CorporationIDs map[int32]struct{}
	AllianceIDs    map[int32]struct{}
}

func (ahs AuthHashSet) ContainsCharacter(id int32) bool {
	_, ok := ahs.CharacterIDs[id]
	return ok
}

func (ahs AuthHashSet) ContainsCorporation(id int32) bool {
	_, ok := ahs.CorporationIDs[id]
	return ok
}

func (ahs AuthHashSet) ContainsAlliance(id int32) bool {
	_, ok := ahs.AllianceIDs[id]
	return ok
}
