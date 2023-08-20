package cachekeys

import "fmt"

func BucketReaderCacheKey(objName string) string {
	return objName
}

func RateLimitingContractItemsCacheKey(contractId int32) string {
	return fmt.Sprintf("contractitems-%d", contractId)
}

func ContractsCacheKey() string {
	return "contracts"
}

func JWKSCacheKey() string {
	return "jwks"
}

func CharacterInfoCacheKey(characterId int32) string {
	return fmt.Sprintf("charinfo-%d", characterId)
}

func StructureInfoCacheKey(structureId int64) string {
	return fmt.Sprintf("structureinfo-%d", structureId)
}

func AllShopAssetsCacheKey() string {
	return "allassets"
}

func UnreservedShopAssetsCacheKey() string {
	return "unresassets"
}

func RegionMarketCacheKey(
	regionId int32,
	typeId int32,
	isBuy bool,
) string {
	return fmt.Sprintf(
		"regionmarket-%d-%d-%s",
		regionId,
		typeId,
		isBuyStr(isBuy),
	)
}

func FilterRegionMarketCacheKey(
	regionId int32,
	typeId int32,
	isBuy bool,
	locationId int64,
) string {
	return fmt.Sprintf(
		"%s-%d",
		RegionMarketCacheKey(regionId, typeId, isBuy),
		locationId,
	)
}

func StructureMarketCacheKey(structureId int64) string {
	return fmt.Sprintf("structuremarket-%d", structureId)
}

func FilterStructureMarketCacheKey(
	structureId int64,
	typeId int32,
	isBuy bool,
) string {
	return fmt.Sprintf(
		"%s-%d-%s",
		StructureMarketCacheKey(structureId),
		typeId,
		isBuyStr(isBuy),
	)
}

func ReadAppraisalCacheKey(appraisalCode string) string {
	return appraisalCode
	// return fmt.Sprintf("appraisal-%s", appraisalCode)
}

func ReadCharacterAppraisalCodesCacheKey(characterId int32) string {
	return fmt.Sprintf("appraisal-codes-%d", characterId)
}

func ShopQueueReadCacheKey() string {
	return "shopqueue"
}

// // util
func isBuyStr(isBuy bool) string {
	if isBuy {
		return "buy"
	} else {
		return "sell"
	}
}
