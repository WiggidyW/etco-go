package tc

import "github.com/WiggidyW/weve-esi/staticdb"

var KVReaderShopInfo tKVReaderShopInfo

type tKVReaderShopInfo struct{}

func (tKVReaderShopInfo) Get(capacity int) *ShopInfo {
	return newShopInfo(capacity)
}

type ShopInfo struct {
	locationMap map[int64]*ShopLocationInfo
}

func newShopInfo(capacity int) *ShopInfo {
	return &ShopInfo{
		locationMap: make(map[int64]*ShopLocationInfo, capacity),
	}
}

func (si *ShopInfo) GetLocation(k int64) (*ShopLocationInfo, bool) {
	shopLocation, ok := si.locationMap[k]
	if !ok {
		if newShopLocation, ok := kVReaderShopLocation.Get(k); ok {
			si.locationMap[k] = newShopLocationInfo(
				newShopLocation,
			)
		} else {
			si.locationMap[k] = nil
		}
	}
	return shopLocation, shopLocation != nil
}

type ShopLocationInfo struct {
	shopLocation ShopLocation
	typeMap      *staticdb.Container[map[int32]int]
	bannedFlags  *staticdb.Container[map[string]struct{}]
}

func newShopLocationInfo(shopLocation ShopLocation) *ShopLocationInfo {
	return &ShopLocationInfo{shopLocation: shopLocation}
}

func (sl *ShopLocationInfo) IsBanned(flag string) bool {
	if sl.bannedFlags == nil {
		if sl.shopLocation.BannedFlagsIndex == nil {
			sl.bannedFlags = staticdb.
				NewContainer[map[string]struct{}](nil)
		} else {
			sl.bannedFlags = staticdb.NewContainer[map[string]struct{}](
				kVReaderBannedFlags.UnsafeGet(
					*sl.shopLocation.BannedFlagsIndex,
				),
			)
		}
	}
	_, ok := sl.bannedFlags.Inner[flag]
	return ok
}

func (sl *ShopLocationInfo) getTypeMap() map[int32]int {
	if sl.typeMap == nil {
		sl.typeMap = staticdb.NewContainer[map[int32]int](
			kVReaderShopTypeMap.UnsafeGet(
				sl.shopLocation.TypeMapIndex,
			),
		)
	}
	return sl.typeMap.Inner
}

func (sl *ShopLocationInfo) GetType(k int32) (*PricingInfo, bool) {
	typeMap := sl.getTypeMap()
	if pricingIndex, ok := typeMap[k]; ok {
		return newPricingInfo(kVReaderPricing.UnsafeGet(
			pricingIndex,
		)), true
	}
	return nil, false
}

func (sl *ShopLocationInfo) HasType(k int32) bool {
	_, ok := sl.getTypeMap()[k]
	return ok
}
