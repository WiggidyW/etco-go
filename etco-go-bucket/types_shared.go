package etcogobucket

type DecPercentage = uint16 // 10,000 = 100.00%, 5,555 = 55.55%, 1 = 0.01%

func NewDecPercentage[U ~uint32 | ~uint64 | ~uint](u U) DecPercentage {
	if u > 10000 {
		return 10000
	}
	return DecPercentage(u)
}

type TypeId = int32
type LocationId = int64
type SystemId = int32
