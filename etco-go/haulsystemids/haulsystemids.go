package haulsystemids

type HaulSystemIds struct {
	Start int32
	End   int32
}

// only used for hashing
func (hsi HaulSystemIds) ToInt64() (i int64) {
	return int64(hsi.Start)<<32 | int64(hsi.End)
}
