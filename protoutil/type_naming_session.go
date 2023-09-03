package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

func MaybeNewLocalTypeNamingSession(
	includeNaming *proto.IncludeTypeNaming,
) *staticdb.TypeNamingSession[*staticdb.LocalIndexMap] {
	if includeNaming == nil {
		return nil
	}
	namingSessionVal := staticdb.NewLocalTypeNamingSession(
		includeNaming.IncludeName,
		includeNaming.IncludeMarketGroups,
		includeNaming.IncludeGroup,
		includeNaming.IncludeCategory,
	)
	return &namingSessionVal
}

func MaybeNewSyncTypeNamingSession(
	includeNaming *proto.IncludeTypeNaming,
) *staticdb.TypeNamingSession[*staticdb.SyncIndexMap] {
	if includeNaming == nil {
		return nil
	}
	namingSessionVal := staticdb.NewSyncTypeNamingSession(
		includeNaming.IncludeName,
		includeNaming.IncludeMarketGroups,
		includeNaming.IncludeGroup,
		includeNaming.IncludeCategory,
	)
	return &namingSessionVal
}

func MaybeGetTypeNamingIndexes[IM staticdb.IndexMap](
	namingSession *staticdb.TypeNamingSession[IM],
	typeId int32,
) *proto.TypeNamingIndexes {
	if namingSession == nil {
		return nil
	}
	rNaming := namingSession.AddType(typeId)
	return newPBTypeNamingIndexes(rNaming)
}

func MaybeFinishTypeNamingSession[IM staticdb.IndexMap](
	namingSession *staticdb.TypeNamingSession[IM],
) *proto.TypeNamingLists {
	if namingSession == nil {
		return nil
	}
	mrktGroups, groups, categories := namingSession.Finish()
	return &proto.TypeNamingLists{
		MarketGroups: mrktGroups,
		Groups:       groups,
		Categories:   categories,
	}
}

func newPBTypeNamingIndexes(
	rTypeNamingIndexes staticdb.TypeNamingIndexes,
) *proto.TypeNamingIndexes {
	return &proto.TypeNamingIndexes{
		Name:               rTypeNamingIndexes.Name,
		GroupIndex:         rTypeNamingIndexes.GroupIndex,
		CategoryIndex:      rTypeNamingIndexes.CategoryIndex,
		MarketGroupIndexes: rTypeNamingIndexes.MarketGroupIndexes,
	}
}
