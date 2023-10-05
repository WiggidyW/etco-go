package staticdb

import "sync"

type TypeNamingIndexes struct {
	Name               string
	MarketGroupIndexes []int32
	GroupIndex         int32
	CategoryIndex      int32
}

func newTypeNamingIndexes() *TypeNamingIndexes {
	return &TypeNamingIndexes{
		Name:               "",
		MarketGroupIndexes: []int32{},
		GroupIndex:         -1,
		CategoryIndex:      -1,
	}
}

type TypeNamingSession[IM IndexMap] struct {
	includeAny          bool
	includeName         bool
	includeMarketGroups bool
	includeGroup        bool
	includeCategory     bool
	marketGroupIM       IM
	groupIM             IM
	categoryIM          IM
}

func NewLocalTypeNamingSession(
	includeName bool,
	includeMarketGroups bool,
	includeGroup bool,
	includeCategory bool,
) TypeNamingSession[*LocalIndexMap] {
	return TypeNamingSession[*LocalIndexMap]{
		includeAny: includeName ||
			includeMarketGroups ||
			includeGroup ||
			includeCategory,
		includeName:         includeName,
		includeMarketGroups: includeMarketGroups,
		includeGroup:        includeGroup,
		includeCategory:     includeCategory,
		marketGroupIM:       newLocalIndexMap(0),
		groupIM:             newLocalIndexMap(0),
		categoryIM:          newLocalIndexMap(0),
	}
}

func NewSyncTypeNamingSession(
	includeName bool,
	includeMarketGroups bool,
	includeGroup bool,
	includeCategory bool,
) TypeNamingSession[*SyncIndexMap] {
	return TypeNamingSession[*SyncIndexMap]{
		includeAny: includeName ||
			includeMarketGroups ||
			includeGroup ||
			includeCategory,
		includeName:         includeName,
		includeMarketGroups: includeMarketGroups,
		includeGroup:        includeGroup,
		includeCategory:     includeCategory,
		marketGroupIM:       newSyncIndexMap(0),
		groupIM:             newSyncIndexMap(0),
		categoryIM:          newSyncIndexMap(0),
	}
}

func (ns TypeNamingSession[IM]) AddType(typeId int32) TypeNamingIndexes {
	n := newTypeNamingIndexes()

	// set indexes only if the type exists and we're including anything
	if ns.includeAny {
		if sdeTypeInfo := GetSDETypeInfo(typeId); sdeTypeInfo != nil {
			if ns.includeName {
				ns.addName(*sdeTypeInfo, n)
			}
			if ns.includeMarketGroups {
				ns.addMarketGroups(*sdeTypeInfo, n)
			}
			if ns.includeGroup {
				ns.addGroup(*sdeTypeInfo, n)
			}
			if ns.includeCategory {
				ns.addCategory(*sdeTypeInfo, n)
			}
		}
	}

	return *n
}

func (ns TypeNamingSession[IM]) Finish() (
	marketGroups,
	groups,
	categories []string,
) {
	marketGroups = ns.marketGroupIM.keys()
	groups = ns.groupIM.keys()
	categories = ns.categoryIM.keys()
	return
}

func (ns TypeNamingSession[IM]) addName(t SDETypeInfo, n *TypeNamingIndexes) {
	n.Name = t.Name
}

func (ns TypeNamingSession[IM]) addMarketGroups(
	t SDETypeInfo,
	n *TypeNamingIndexes,
) {
	// continue only if the type has any market groups
	if t.MarketGroups != nil && len(t.MarketGroups) > 0 {
		// initialize the index slice
		n.MarketGroupIndexes = make([]int32, 0, len(t.MarketGroups))
		// add the index of each market group to the slice
		for _, marketGroup := range t.MarketGroups {
			// try to get existing index
			idx, ok := ns.marketGroupIM.get(marketGroup)
			// otherwise, add the new string
			if !ok {
				idx = ns.marketGroupIM.add(marketGroup)
			}
			// add the index to the list
			n.MarketGroupIndexes = append(
				n.MarketGroupIndexes,
				idx,
			)
		}
	}
}

func (ns TypeNamingSession[IM]) addGroup(t SDETypeInfo, n *TypeNamingIndexes) {
	// try to get existing index
	idx, ok := ns.groupIM.get(t.Group)
	// otherwise, add the new string
	if !ok {
		idx = ns.groupIM.add(t.Group)
	}
	// set the index
	n.GroupIndex = idx
}

func (ns TypeNamingSession[IM]) addCategory(
	t SDETypeInfo,
	n *TypeNamingIndexes,
) {
	// try to get existing index
	idx, ok := ns.categoryIM.get(t.Category)
	// otherwise, add the new string
	if !ok {
		idx = ns.categoryIM.add(t.Category)
	}
	// set the index
	n.CategoryIndex = idx
}

type IndexMap interface {
	keys() []string
	get(key string) (index int32, ok bool)
	add(key string) (index int32)
}

type LocalIndexMap struct {
	keys_      []string
	keyIndexes map[string]int32
}

func newLocalIndexMap(capacity int) *LocalIndexMap {
	return &LocalIndexMap{
		keys_:      make([]string, 0, capacity),
		keyIndexes: make(map[string]int32, capacity),
	}
}

func (bim *LocalIndexMap) keys() []string {
	return bim.keys_
}

func (bim *LocalIndexMap) get(key string) (index int32, ok bool) {
	index, ok = bim.keyIndexes[key]
	return index, ok
}

func (bim *LocalIndexMap) add(key string) (index int32) {
	index = int32(len(bim.keys_))
	bim.keyIndexes[key] = index
	bim.keys_ = append(bim.keys_, key)
	return index
}

type SyncIndexMap struct {
	keys_      []string
	keyIndexes map[string]int32
	rwLock     *sync.RWMutex
}

func newSyncIndexMap(capacity int) *SyncIndexMap {
	return &SyncIndexMap{
		keys_:      make([]string, 0, capacity),
		keyIndexes: make(map[string]int32, capacity),
		rwLock:     &sync.RWMutex{},
	}
}

func (sim *SyncIndexMap) keys() []string {
	sim.rwLock.RLock()
	defer sim.rwLock.RUnlock()
	return sim.keys_
}

func (sim *SyncIndexMap) get(key string) (index int32, ok bool) {
	sim.rwLock.RLock()
	defer sim.rwLock.RUnlock()
	index, ok = sim.keyIndexes[key]
	return index, ok
}

func (sim *SyncIndexMap) add(key string) (index int32) {
	sim.rwLock.Lock()
	defer sim.rwLock.Unlock()
	index = int32(len(sim.keys_))
	sim.keyIndexes[key] = index
	sim.keys_ = append(sim.keys_, key)
	return index
}
