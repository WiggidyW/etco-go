package sde

import "sync"

type Naming struct {
	name               string
	marketGroupIndexes []int
	groupIndex         int
	categoryIndex      int
}

func newNaming() *Naming {
	return &Naming{
		name:               "",
		marketGroupIndexes: []int{},
		groupIndex:         -1,
		categoryIndex:      -1,
	}
}

type NamingSession struct {
	includeName         bool
	includeMarketGroups bool
	includeGroup        bool
	includeCategory     bool
	marketGroupIM       syncIndexMap
	groupIM             syncIndexMap
	categoryIM          syncIndexMap
}

func (ns *NamingSession) addName(t *TypeInfo, n *Naming) {
	n.name = t.Name()
}

func (ns *NamingSession) addMarketGroups(t *TypeInfo, n *Naming) {
	// continue only if the type has any market groups
	if marketGroups, ok := t.MarketGroups(); ok {
		// initialize the index slice
		n.marketGroupIndexes = make([]int, len(marketGroups))
		// add the index of each market group to the slice
		for _, marketGroup := range marketGroups {
			// try to get existing index
			idx, ok := ns.marketGroupIM.Get(marketGroup)
			// otherwise, add the new string
			if !ok {
				idx = ns.marketGroupIM.Add(marketGroup)
			}
			// add the index to the list
			n.marketGroupIndexes = append(
				n.marketGroupIndexes,
				idx,
			)
		}
	}
}

func (ns *NamingSession) addGroup(t *TypeInfo, n *Naming) {
	// continue only if the type has a group
	if group, ok := t.Group(); ok {
		// try to get existing index
		idx, ok := ns.groupIM.Get(group)
		// otherwise, add the new string
		if !ok {
			idx = ns.groupIM.Add(group)
		}
		// set the index
		n.groupIndex = idx
	}
}

func (ns *NamingSession) addCategory(t *TypeInfo, n *Naming) {
	// continue only if the type has a category
	if category, ok := t.Category(); ok {
		// try to get existing index
		idx, ok := ns.categoryIM.Get(category)
		// otherwise, add the new string
		if !ok {
			idx = ns.categoryIM.Add(category)
		}
		// set the index
		n.categoryIndex = idx
	}
}

func (ns *NamingSession) AddType(typeId int32) Naming {
	n := newNaming()
	typeInfo, ok := KVReaderTypeInfo.Get(typeId)
	if ok { // set indexes only if the type exists
		if ns.includeName {
			ns.addName(typeInfo, n)
		}
		if ns.includeMarketGroups {
			ns.addMarketGroups(typeInfo, n)
		}
		if ns.includeGroup {
			ns.addGroup(typeInfo, n)
		}
		if ns.includeCategory {
			ns.addCategory(typeInfo, n)
		}
	}
	return *n
}

func (ns *NamingSession) Finish() (marketGroups, groups, categories []string) {
	marketGroups = ns.marketGroupIM.keys()
	groups = ns.groupIM.keys()
	categories = ns.categoryIM.keys()
	return
}

type syncIndexMap struct {
	keys_      []string
	keyIndexes map[string]int
	rwLock     *sync.RWMutex
}

func (sim *syncIndexMap) keys() []string {
	sim.rwLock.RLock()
	defer sim.rwLock.RUnlock()
	return sim.keys_
}

func (sim *syncIndexMap) Get(key string) (index int, ok bool) {
	sim.rwLock.RLock()
	defer sim.rwLock.RUnlock()
	index, ok = sim.keyIndexes[key]
	return
}

func (sim *syncIndexMap) Add(key string) (index int) {
	sim.rwLock.Lock()
	defer sim.rwLock.Unlock()
	index = len(sim.keys_)
	sim.keyIndexes[key] = index
	sim.keys_ = append(sim.keys_, key)
	return
}
