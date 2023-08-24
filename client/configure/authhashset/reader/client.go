package reader

import (
	"github.com/WiggidyW/eve-trading-co-go/client/configure/authhashset"
	b "github.com/WiggidyW/eve-trading-co-go/client/configure/internal/bucket/reader"
)

type AuthHashSetReaderClient = b.SC_BucketReaderClient[authhashset.AuthHashSet]
