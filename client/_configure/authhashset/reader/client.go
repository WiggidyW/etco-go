package reader

import (
	b "github.com/WiggidyW/etco-go-bucket"

	breader "github.com/WiggidyW/etco-go/client/configure/internal/bucket/reader"
)

type AuthHashSetReaderClient = breader.SC_BucketReaderClient[b.AuthHashSet]
