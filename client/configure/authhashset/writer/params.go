package writer

import (
	ahs "github.com/WiggidyW/weve-esi/client/configure/authhashset"
	b "github.com/WiggidyW/weve-esi/client/configure/internal/bucket/writer"
)

type AuthHashSetWriterParams = b.BucketWriterParams[ahs.AuthHashSet]
