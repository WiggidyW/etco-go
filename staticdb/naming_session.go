package staticdb

import "sync"

const (
	MARKET_GROUPS_PER_TYPE int = 5
	GROUPS_PER_TYPE        int = 1
	CATEGORIES_PER_TYPE    int = 1
)

type Naming struct {
	name             string
	mrktGroupIndexes []int
	groupIndex       int
	categoryIndex    int
}

func newNaming() *Naming {
	return &Naming{
		name:             "",
		mrktGroupIndexes: []int{},
		groupIndex:       -1,
		categoryIndex:    -1,
	}
}

type NamingSession struct {
	includeAny        bool
	includeName       bool
	includeMrktGroups bool
	includeGroup      bool
	includeCategory   bool
	mrktGroupIM       *syncIndexMap
	groupIM           *syncIndexMap
	categoryIM        *syncIndexMap
}

func NewNamingSession(
	includeName bool,
	includeMrktGroups bool,
	includeGroup bool,
	includeCategory bool,
) NamingSession {
	return NamingSession{
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

func (ns NamingSession) AddType(typeId int32) Naming {
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

func (ns NamingSession) Finish() (mrktGroups, groups, categories []string) {
	mrktGroups = ns.mrktGroupIM.keys()
	groups = ns.groupIM.keys()
	categories = ns.categoryIM.keys()
	return
}

func (ns NamingSession) addName(t SDETypeInfo, n *Naming) {
	n.name = t.Name
}

func (ns NamingSession) addMrktGroups(t SDETypeInfo, n *Naming) {
	// continue only if the type has any mrkt groups
	if t.MrktGroups != nil && len(t.MrktGroups) > 0 {
		// initialize the index slice
		n.mrktGroupIndexes = make([]int, len(t.MrktGroups))
		// add the index of each mrkt group to the slice
		for _, mrktGroup := range t.MrktGroups {
			// try to get existing index
			idx, ok := ns.mrktGroupIM.get(mrktGroup)
			// otherwise, add the new string
			if !ok {
				idx = ns.mrktGroupIM.add(mrktGroup)
			}
			// add the index to the list
			n.mrktGroupIndexes = append(
				n.mrktGroupIndexes,
				idx,
			)
		}
	}
}

func (ns NamingSession) addGroup(t SDETypeInfo, n *Naming) {
	// continue only if the type has a group
	if t.Group != nil {
		// try to get existing index
		idx, ok := ns.groupIM.get(*t.Group)
		// otherwise, add the new string
		if !ok {
			idx = ns.groupIM.add(*t.Group)
		}
		// set the index
		n.groupIndex = idx
	}
}

func (ns NamingSession) addCategory(t SDETypeInfo, n *Naming) {
	// continue only if the type has a category
	if t.Category != nil {
		// try to get existing index
		idx, ok := ns.categoryIM.get(*t.Category)
		// otherwise, add the new string
		if !ok {
			idx = ns.categoryIM.add(*t.Category)
		}
		// set the index
		n.categoryIndex = idx
	}
}

type syncIndexMap struct {
	keys_      []string
	keyIndexes map[string]int
	rwLock     *sync.RWMutex
}

func newSyncIndexMap(capacity int) *syncIndexMap {
	return &syncIndexMap{
		keys_:      make([]string, 0, capacity),
		keyIndexes: make(map[string]int, capacity),
		rwLock:     &sync.RWMutex{},
	}
}

func (sim *syncIndexMap) keys() []string {
	sim.rwLock.RLock()
	defer sim.rwLock.RUnlock()
	return sim.keys_
}

func (sim *syncIndexMap) get(key string) (index int, ok bool) {
	sim.rwLock.RLock()
	defer sim.rwLock.RUnlock()
	index, ok = sim.keyIndexes[key]
	return
}

func (sim *syncIndexMap) add(key string) (index int) {
	sim.rwLock.Lock()
	defer sim.rwLock.Unlock()
	index = len(sim.keys_)
	sim.keyIndexes[key] = index
	sim.keys_ = append(sim.keys_, key)
	return
}
