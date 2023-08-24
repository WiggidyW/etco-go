package staticdb

import "sync"

// const (
// 	MARKET_GROUPS_PER_TYPE int = 5
// 	GROUPS_PER_TYPE        int = 1
// 	CATEGORIES_PER_TYPE    int = 1
// )

type Naming struct {
	Name             string
	MrktGroupIndexes []int32
	GroupIndex       int32
	CategoryIndex    int32
}

func newNaming() *Naming {
	return &Naming{
		Name:             "",
		MrktGroupIndexes: []int32{},
		GroupIndex:       -1,
		CategoryIndex:    -1,
	}
}

type NamingSession[IM IndexMap] struct {
	includeAny        bool
	includeName       bool
	includeMrktGroups bool
	includeGroup      bool
	includeCategory   bool
	mrktGroupIM       IM
	groupIM           IM
	categoryIM        IM
}

func NewLocalNamingSession(
	includeName bool,
	includeMrktGroups bool,
	includeGroup bool,
	includeCategory bool,
) NamingSession[*LocalIndexMap] {
	return NamingSession[*LocalIndexMap]{
		includeAny: includeName ||
			includeMrktGroups ||
			includeGroup ||
			includeCategory,
		includeName:       includeName,
		includeMrktGroups: includeMrktGroups,
		includeGroup:      includeGroup,
		includeCategory:   includeCategory,
		mrktGroupIM:       newLocalIndexMap(0),
		groupIM:           newLocalIndexMap(0),
		categoryIM:        newLocalIndexMap(0),
	}
}

func NewSyncNamingSession(
	includeName bool,
	includeMrktGroups bool,
	includeGroup bool,
	includeCategory bool,
) NamingSession[*SyncIndexMap] {
	return NamingSession[*SyncIndexMap]{
		includeAny: includeName ||
			includeMrktGroups ||
			includeGroup ||
			includeCategory,
		includeName:       includeName,
		includeMrktGroups: includeMrktGroups,
		includeGroup:      includeGroup,
		includeCategory:   includeCategory,
		mrktGroupIM:       newSyncIndexMap(0),
		groupIM:           newSyncIndexMap(0),
		categoryIM:        newSyncIndexMap(0),
	}
}

func (ns NamingSession[IM]) AddType(typeId int32) Naming {
	n := newNaming()

	// set indexes only if the type exists and we're including anything
	if ns.includeAny {
		if typeInfo := GetSDETypeInfo(typeId); typeInfo != nil {
			if ns.includeName {
				ns.addName(*typeInfo, n)
			}
			if ns.includeMrktGroups {
				ns.addMrktGroups(*typeInfo, n)
			}
			if ns.includeGroup {
				ns.addGroup(*typeInfo, n)
			}
			if ns.includeCategory {
				ns.addCategory(*typeInfo, n)
			}
		}
	}

	return *n
}

func (ns NamingSession[IM]) Finish() (mrktGroups, groups, categories []string) {
	mrktGroups = ns.mrktGroupIM.keys()
	groups = ns.groupIM.keys()
	categories = ns.categoryIM.keys()
	return
}

func (ns NamingSession[IM]) addName(t SDETypeInfo, n *Naming) {
	n.Name = t.Name
}

func (ns NamingSession[IM]) addMrktGroups(t SDETypeInfo, n *Naming) {
	// continue only if the type has any mrkt groups
	if t.MrktGroups != nil && len(t.MrktGroups) > 0 {
		// initialize the index slice
		n.MrktGroupIndexes = make([]int32, len(t.MrktGroups))
		// add the index of each mrkt group to the slice
		for _, mrktGroup := range t.MrktGroups {
			// try to get existing index
			idx, ok := ns.mrktGroupIM.get(mrktGroup)
			// otherwise, add the new string
			if !ok {
				idx = ns.mrktGroupIM.add(mrktGroup)
			}
			// add the index to the list
			n.MrktGroupIndexes = append(
				n.MrktGroupIndexes,
				idx,
			)
		}
	}
}

func (ns NamingSession[IM]) addGroup(t SDETypeInfo, n *Naming) {
	// continue only if the type has a group
	if t.Group != nil {
		// try to get existing index
		idx, ok := ns.groupIM.get(*t.Group)
		// otherwise, add the new string
		if !ok {
			idx = ns.groupIM.add(*t.Group)
		}
		// set the index
		n.GroupIndex = idx
	}
}

func (ns NamingSession[IM]) addCategory(t SDETypeInfo, n *Naming) {
	// continue only if the type has a category
	if t.Category != nil {
		// try to get existing index
		idx, ok := ns.categoryIM.get(*t.Category)
		// otherwise, add the new string
		if !ok {
			idx = ns.categoryIM.add(*t.Category)
		}
		// set the index
		n.CategoryIndex = idx
	}
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
