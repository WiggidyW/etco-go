package etcogobucket

type BucketKind uint8

const (
	WEB BucketKind = iota
	AUTH
	BUILD
)
