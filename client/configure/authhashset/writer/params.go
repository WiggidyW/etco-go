package writer

import (
	ahs "github.com/WiggidyW/eve-trading-co-go/client/configure/authhashset"
	b "github.com/WiggidyW/eve-trading-co-go/client/configure/internal/bucket/writer"
)

type AuthHashSetWriterParams = b.BucketWriterParams[ahs.AuthHashSet]
