package assets

import (
	"github.com/WiggidyW/weve-esi/client/modelclient"
	"github.com/WiggidyW/weve-esi/staticdb"
)

type ShopAsset struct {
	TypeId int32
	// upgrade quantity to int64 because we'll be summing int32s
	Quantity int64
}

type flattenedAsset struct {
	Flags      []string
	LocationId int64
	ShopAsset  ShopAsset
}

type unflattenedAsset struct {
	ItemId     int64
	LocationId int64
	Quantity   int32
	TypeId     int32
}

type locationAndFlag struct {
	LocationId int64
	Flag       string
}

type unflattenedAssets struct {
	locationAndFlags map[int64] /*itemId*/ locationAndFlag
	assets           []unflattenedAsset
	flagBuf          *[]string
	// Deduplicate
	// map[locationid]map[typeid]*ShopAsset -> map[locationid][]ShopAsset
	shopAssetsDeduper map[int64]map[int32]*ShopAsset
}

func newUnflattenedAssets(capacity int32) *unflattenedAssets {
	return &unflattenedAssets{
		locationAndFlags:  make(map[int64]locationAndFlag, capacity),
		assets:            make([]unflattenedAsset, 0, capacity),
		flagBuf:           new([]string),
		shopAssetsDeduper: make(map[int64]map[int32]*ShopAsset),
	}
}

func (ua *unflattenedAssets) addEntry(entry modelclient.EntryAssetsCorporation) {
	ua.locationAndFlags[entry.ItemId] = locationAndFlag{
		LocationId: entry.LocationId,
		Flag:       entry.LocationFlag,
	}
	ua.assets = append(ua.assets, unflattenedAsset{
		ItemId:     entry.ItemId,
		LocationId: entry.LocationId,
		Quantity:   entry.Quantity,
		TypeId:     entry.TypeId,
	})
}

// freezes the current flag buffer, returns it and clears it
func (ua *unflattenedAssets) currentFlags() []string {
	currentFlags := *ua.flagBuf
	*ua.flagBuf = (*ua.flagBuf)[:0]
	return currentFlags
}

// flattens an asset recursively by traversing locationAndFlags
func (ua *unflattenedAssets) flattenAsset(a unflattenedAsset) flattenedAsset {
	laf := ua.locationAndFlags[a.ItemId]
	for {
		*ua.flagBuf = append(*ua.flagBuf, laf.Flag)
		parentLaf, ok := ua.locationAndFlags[laf.LocationId]
		if ok { // if a parent exists, continue
			laf = parentLaf
		} else { // else, break
			break
		}
	}
	return flattenedAsset{
		Flags:      ua.currentFlags(),
		LocationId: laf.LocationId,
		ShopAsset: ShopAsset{
			Quantity: int64(a.Quantity),
			TypeId:   a.TypeId,
		},
	}
}

// checks if an asset should be filtered out
func (ua *unflattenedAssets) okFlattenedAsset(a flattenedAsset) bool {
	// if location is not in shopInfo, return false
	shopInfo := staticdb.GetShopLocationInfo(a.LocationId)
	if shopInfo == nil {
		return false
	}
	// if type is not in locationInfo, return false
	if !shopInfo.HasTypeInfo(a.ShopAsset.TypeId) {
		return false
	}
	// if any flag is banned, return false
	for _, flag := range a.Flags {
		if shopInfo.BannedFlags.Has(flag) {
			return false
		}
	}
	return true
}

// appends a flattened asset to the shopAssetDeduper
func (ua *unflattenedAssets) addFlattenedAsset(a flattenedAsset) {
	locationMap, ok := ua.shopAssetsDeduper[a.LocationId]
	if !ok {
		locationMap = make(map[int32]*ShopAsset)
		ua.shopAssetsDeduper[a.LocationId] = locationMap
	}
	if shopAsset, ok := locationMap[a.ShopAsset.TypeId]; ok {
		shopAsset.Quantity += a.ShopAsset.Quantity
	} else {
		locationMap[a.ShopAsset.TypeId] = &a.ShopAsset
	}
}

// converts ua.shopAssetsDeduper - map[int64]map[int32]*ShopAsset
// to the output - map[int64][]ShopAsset
func (ua *unflattenedAssets) typeMapsToSlices() map[int64][]ShopAsset {
	shopAssets := make(map[int64][]ShopAsset, len(ua.shopAssetsDeduper))
	for locationId, typeMap := range ua.shopAssetsDeduper {
		locationSlice := make([]ShopAsset, 0, len(typeMap))
		for _, shopAsset := range typeMap {
			locationSlice = append(locationSlice, *shopAsset)
		}
		shopAssets[locationId] = locationSlice
	}
	return shopAssets
}

// flattens and filters the inner assets, converting them to valid ShopAssets
func (ua *unflattenedAssets) flattenAndFilter() map[int64][]ShopAsset {
	// flatten and filter ua.assets
	for _, unfltAsset := range ua.assets {
		// flatten the asset
		fltAsset := ua.flattenAsset(unfltAsset)
		// filter the asset
		if !ua.okFlattenedAsset(fltAsset) {
			continue
		}
		// append the asset
		ua.addFlattenedAsset(fltAsset)
	}
	// convert ua.shopAssetsDeduper to the output
	return ua.typeMapsToSlices()
}
