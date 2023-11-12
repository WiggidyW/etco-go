package appraisal

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"sync"
	"time"
)

var (
	hasherPool sync.Pool = sync.Pool{New: func() any { return fnv.New64() }}
)

func hashAppraisal[TERID ~int32 | ~int64](
	codeChar byte,
	time time.Time,
	itemsLength int,
	version string,
	characterIdPtr *int32,
	territoryId TERID,
	price, tax, taxRate, fee, feePerM3 float64,
) (code string) {
	hasher := hasherPool.Get().(hash.Hash64)

	var characterId int32 = 0
	if characterIdPtr != nil {
		characterId = *characterIdPtr
	}

	hasher.Write(i64Bytes(time.Unix()))
	hasher.Write(i64Bytes(int64(itemsLength)))
	hasher.Write([]byte(version))
	hasher.Write(i32Bytes(characterId))
	hasher.Write(i64Bytes(int64(territoryId)))
	hasher.Write(f64Bytes(price))
	hasher.Write(f64Bytes(tax))
	hasher.Write(f64Bytes(taxRate))
	hasher.Write(f64Bytes(fee))
	hasher.Write(f64Bytes(feePerM3))

	code = string(codeChar) + fmt.Sprintf("%016x", hasher.Sum64())[1:]

	hasher.Reset()
	hasherPool.Put(hasher)

	return code
}

// func getHasher() hash.Hash64 {
// 	return hasherPool.Get().(hash.Hash64)
// }

// func hashItem[AITEM AppraisalItem](
// 	hasher hash.Hash64,
// 	item AITEM,
// ) {
// 	hasher.Write(i32Bytes(item.GetTypeId()))
// 	hasher.Write(i64Bytes(item.GetQuantity()))
// 	hasher.Write(f64Bytes(item.GetPricePerUnit()))
// 	hasher.Write([]byte(item.GetDescription()))
// 	hasher.Write(f64Bytes(item.GetFeePerUnit()))
// 	hasher.Write(i64Bytes(int64(item.GetChildrenLength())))
// }

// func hashFinish(
// 	hasher hash.Hash64,
// 	codeChar byte,
// ) string {
// 	code := string(codeChar) + fmt.Sprintf("%016x", hasher.Sum64())[1:]
// 	hasher.Reset()
// 	hasherPool.Put(hasher)
// 	return code
// }

func i32Bytes[N ~int32 | ~uint32](i N) []byte {
	return []byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

func i64Bytes[N ~int64 | ~uint64](i N) []byte {
	return []byte{
		byte(i >> 56),
		byte(i >> 48),
		byte(i >> 40),
		byte(i >> 32),
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

func f64Bytes(f float64) []byte {
	return i64Bytes(math.Float64bits(f))
}
