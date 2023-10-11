package allassets_

import (
	modelac "github.com/WiggidyW/etco-go/client/esi/model/assetscorporation"
	"github.com/WiggidyW/etco-go/staticdb"
)

type flattenedAsset struct {
	Flags      []string
	LocationId int64
	TypeId     int32
	Quantity   int64
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

type shopLocationMap struct {
	inner map[int64]*staticdb.ShopLocationInfo
}

func (slm shopLocationMap) getShopLocationInfo(
	k int64,
) *staticdb.ShopLocationInfo {
	sli, ok := slm.inner[k]
	if !ok {
		sli = staticdb.GetShopLocationInfo(k)
		slm.inner[k] = sli
	}
	return sli
}

// TODO: Break this up into smaller functions / structs
type unflattenedAssets struct {
	locationAndFlags map[int64] /*itemId*/ locationAndFlag
	assets           []unflattenedAsset
	flagBuf          *[]string
	shopLocationMap  shopLocationMap
	shopAssets       map[int64]map[int32]*int64
}

func newUnflattenedAssets(capacity int) *unflattenedAssets {
	return &unflattenedAssets{
		locationAndFlags: make(map[int64]locationAndFlag, capacity),
		assets:           make([]unflattenedAsset, 0, capacity),
		flagBuf:          &[]string{},
		shopLocationMap: shopLocationMap{
			inner: make(map[int64]*staticdb.ShopLocationInfo),
		},
		shopAssets: make(map[int64]map[int32]*int64),
	}
}

func (ua *unflattenedAssets) addEntry(entry modelac.AssetsCorporationEntry) {
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
		TypeId:     a.TypeId,
		Quantity:   int64(a.Quantity),
	}
}

// checks if an asset should be filtered out
func (ua *unflattenedAssets) okFlattenedAsset(a flattenedAsset) bool {
	// if location is not in shopInfo, return false
	shopInfo := ua.shopLocationMap.getShopLocationInfo(a.LocationId)
	if shopInfo == nil {
		return false
	}
	// if type is not in locationInfo, return false
	if !shopInfo.HasTypePricingInfo(a.TypeId) {
		return false
	}
	// if any flag is banned, return false
	for _, flag := range a.Flags {
		if shopInfo.HasBannedFlag(flag) {
			return false
		}
	}
	return true
}

// appends a flattened asset to the shopAssetDeduper
func (ua *unflattenedAssets) addFlattenedAsset(fltAsset flattenedAsset) {
	// get the location map or create if not exists
	locationMap, ok := ua.shopAssets[fltAsset.LocationId]
	if !ok {
		locationMap = make(map[int32]*int64)
		ua.shopAssets[fltAsset.LocationId] = locationMap
	}
	// if it exists, increase quantity of the existing shopasset
	// else, set the typeid to the new shopasset
	if shopAssetQuant, ok := locationMap[fltAsset.TypeId]; ok {
		*shopAssetQuant += fltAsset.Quantity
	} else {
		newShopAssetQuant := fltAsset.Quantity // allow 'fltAsset' to be GC'd
		locationMap[fltAsset.TypeId] = &newShopAssetQuant
	}
}

// flattens and filters the inner assets, converting them to valid inventory.ShopAssets
func (
	ua *unflattenedAssets,
) flattenAndFilter() map[int64]map[int32]*int64 {
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
	return ua.shopAssets
}
