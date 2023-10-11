package bucket

const (
	WEB_CAPACITY_MULTIPLIER = 3
	WEB_CAPACITY_DIVISOR    = 2
)

func transformWebCapacity(capacity int) int {
	return capacity * WEB_CAPACITY_MULTIPLIER / WEB_CAPACITY_DIVISOR
}
