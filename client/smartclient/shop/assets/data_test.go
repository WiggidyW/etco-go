package assets

// import (
// 	"math/rand"
// 	"testing"

// 	"github.com/WiggidyW/weve-esi/client/modelclient"
// )

// func Benchmark_flatten(b *testing.B) {
// 	var L1Len int = 20000
// 	var L2Len int = 30000
// 	var L3Len int = 40000
// 	var L4Len int = 30000
// 	var L5Len int = 20000
// 	var ParentLocationsLen int = 10

// 	var ParentLocations []int64
// 	for i := 0; i < ParentLocationsLen; i++ {
// 		ParentLocations = append(ParentLocations, rand.Int63())
// 	}
// 	var Flags = []string{"flag1", "flag2", "flag3", "flag4", "flag5"}

// 	entries := make(
// 		[]modelclient.EntryAssetsCorporation,
// 		0,
// 		L1Len+L2Len+L3Len+L4Len+L5Len,
// 	)
// 	var i int = 0

// 	addChildEntries := func(sumLen, prevLen, thisLen int) {
// 		for ; i < sumLen+thisLen; i++ {
// 			entries = append(
// 				entries,
// 				modelclient.EntryAssetsCorporation{
// 					ItemId: rand.Int63(),
// 					LocationId: entries[rand.Intn(
// 						prevLen,
// 					)+sumLen-prevLen].ItemId,
// 					LocationFlag: Flags[rand.Intn(
// 						len(Flags),
// 					)],
// 					Quantity: rand.Int31(),
// 					TypeId:   rand.Int31(),
// 				},
// 			)
// 		}
// 	}

// 	// top level entries
// 	for ; i < L1Len; i++ {
// 		entries = append(
// 			entries,
// 			modelclient.EntryAssetsCorporation{
// 				ItemId: rand.Int63(),
// 				LocationId: ParentLocations[rand.Intn(
// 					ParentLocationsLen,
// 				)],
// 				LocationFlag: Flags[rand.Intn(len(Flags))],
// 				Quantity:     rand.Int31(),
// 				TypeId:       rand.Int31(),
// 			},
// 		)
// 	}
// 	// child entries
// 	addChildEntries(L1Len, L1Len, L2Len)
// 	addChildEntries(L1Len+L2Len, L2Len, L3Len)
// 	addChildEntries(L1Len+L2Len+L3Len, L3Len, L4Len)
// 	addChildEntries(L1Len+L2Len+L3Len+L4Len, L4Len, L5Len)

// 	unflattenedAssets := newUnflattenedAssets(
// 		L1Len + L2Len + L3Len + L4Len + L5Len,
// 	)

// 	b.ResetTimer()
// 	for _, entry := range entries {
// 		unflattenedAssets.addEntry(entry)
// 	}
// 	assets := unflattenedAssets.flattenFilterAssets()

// 	firstLen := len(assets[0].Flags)
// 	success := false
// 	for _, asset := range assets {
// 		if len(asset.Flags) != firstLen {
// 			success = true
// 			break
// 		}
// 	}
// 	if !success {
// 		b.Error("All flattened assets have the same number of flags")
// 	}
// }
