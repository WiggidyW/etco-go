package sde

import "github.com/WiggidyW/weve-esi/staticdb"

var KVReaderTypeInfo tKVReaderTypeInfo

type tKVReaderTypeInfo struct{}

func (tKVReaderTypeInfo) Get(k int32) (v *TypeInfo, ok bool) {
	if typeData, ok := kVReaderTypeData.Get(k); ok {
		return newTypeInfo(typeData), true
	}
	return nil, false
}

func (tKVReaderTypeInfo) UnsafeGet(k int32) *TypeInfo {
	typeData, _ := kVReaderTypeData.Get(k)
	return newTypeInfo(typeData)
}

type TypeInfo struct {
	typeData     TypeData
	marketGroups *staticdb.Container[[]string]
	group        *staticdb.Container[*Group]
	category     *staticdb.Container[*string]
	volume       *staticdb.Container[*float64]
}

func (t *TypeInfo) Name() string {
	return t.typeData.Name
}

func (t *TypeInfo) ReprocessedMaterials() []ReprocessedMaterial {
	return t.typeData.ReprocessedMaterials
}

func (t *TypeInfo) Volume() (float64, bool) {
	if t.volume == nil {
		if t.typeData.VolumeIndex != nil {
			v := kVReaderVolumes.UnsafeGet(*t.typeData.VolumeIndex)
			t.volume = staticdb.NewContainer[*float64](&v)
		} else {
			t.volume = staticdb.NewContainer[*float64](nil)
		}
	}
	if t.volume.Inner == nil {
		return 0, false
	}
	return *t.volume.Inner, true
}

func (t *TypeInfo) Group() (string, bool) {
	group := t.getGroup()
	if group == nil {
		return "", false
	}
	return group.Name, true
}

func (t *TypeInfo) Category() (string, bool) {
	if t.category == nil {
		group := t.getGroup()
		if group != nil && group.CategoryIndex != nil {
			c := kVReaderCategories.UnsafeGet(*group.CategoryIndex)
			t.category = staticdb.NewContainer[*string](&c)
		} else {
			t.category = staticdb.NewContainer[*string](nil)
		}
	}
	if t.category.Inner == nil {
		return "", false
	}
	return *t.category.Inner, true
}

func (t *TypeInfo) MarketGroups() ([]string, bool) {
	if t.marketGroups == nil {
		if t.typeData.MarketGroupIndex != nil {
			// append market groups, starting with inner-most
			var mgidx *int = t.typeData.MarketGroupIndex
			names := make([]string, 0, 1)
			for mgidx != nil { // nil means we've reached the root
				mg := kVReaderMarketGroups.UnsafeGet(*mgidx)
				names = append(names, mg.Name)
				mgidx = mg.ParentIndex
			}
			t.marketGroups = staticdb.NewContainer[[]string](names)
		} else {
			t.marketGroups = staticdb.NewContainer[[]string](nil)
		}
	}
	return t.marketGroups.Inner, t.marketGroups.Inner != nil
}

func (t *TypeInfo) getGroup() *Group {
	if t.group == nil {
		if t.typeData.GroupIndex != nil {
			group := kVReaderGroups.UnsafeGet(
				*t.typeData.GroupIndex,
			)
			t.group = staticdb.NewContainer[*Group](&group)
		} else {
			t.group = staticdb.NewContainer[*Group](nil)
		}
	}
	return t.group.Inner
}

func newTypeInfo(typeData TypeData) *TypeInfo {
	return &TypeInfo{typeData: typeData}
}
