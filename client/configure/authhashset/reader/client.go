package reader

import (
	"github.com/WiggidyW/weve-esi/client/configure/authhashset"
	b "github.com/WiggidyW/weve-esi/client/configure/internal/bucket/reader"
)

type AuthHashSetReaderClient = b.SC_BucketReaderClient[authhashset.AuthHashSet]
