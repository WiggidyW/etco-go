package etcogobucket

import (
	"math"
)

type Scientific16 struct {
	Base   uint8 // 0 - 255 (any)
	Zeroes uint8
}

var int64Pow10Tab = [...]int64{
	1, // 0 zeroes
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000, // 10 zeroes
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,   // 16 zeroes
	100000000000000000,  // 17 zeroes
	1000000000000000000, // 18 zeroes
	// 9223372036854775807,
}
var uint64Pow10Tab = [...]uint64{
	1, // 0 zeroes
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000, // 10 zeroes
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,    // 16 zeroes
	100000000000000000,   // 17 zeroes
	1000000000000000000,  // 18 zeroes
	10000000000000000000, // 19 zeroes
	// 18446744073709551615,
}
var uint32Pow10Tab = [...]uint32{
	1, // 0 zeroes
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,  // 8 zeroes
	1000000000, // 9 zeroes
	4294967295,
}

// may return max int64 if overflow
func (s16 Scientific16) Int64() int64 {
	if s16.Zeroes < 17 ||
		(s16.Zeroes == 17 && s16.Base < 93) ||
		(s16.Zeroes == 18 && s16.Base < 10) {
		return int64(s16.Base) * int64Pow10Tab[s16.Zeroes]
	} else {
		return math.MaxInt64
	}
}

// may return max uint64 if overflow
func (s16 Scientific16) Uint64() uint64 {
	if s16.Zeroes < 17 ||
		(s16.Zeroes == 17 && s16.Base < 185) ||
		(s16.Zeroes == 18 && s16.Base < 19) ||
		(s16.Zeroes == 19 && s16.Base == 1 /* lol */) {
		return uint64(s16.Base) * uint64Pow10Tab[s16.Zeroes]
	} else {
		return math.MaxUint64
	}
}

// may return max uint32 if overflow
func (s16 Scientific16) Uint32() uint32 {
	if s16.Zeroes < 8 ||
		(s16.Zeroes == 8 && s16.Base < 43) ||
		(s16.Zeroes == 9 && s16.Base < 5) {
		return uint32(s16.Base) * uint32Pow10Tab[s16.Zeroes]
	} else {
		return math.MaxUint32
	}
}

func (s16 Scientific16) Float64() float64 {
	return float64(s16.Base) * math.Pow10(int(s16.Zeroes))
}

func NewScientific16FromUInt[U ~uint32 | ~uint64 | ~uint](u U) (
	s16 Scientific16,
) {
	for u > 255 {
		u /= 10
		s16.Zeroes++
	}
	s16.Base = uint8(u)
	return s16
}
