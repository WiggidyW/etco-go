package writer

import (
	ahs "github.com/WiggidyW/weve-esi/client/configure/authhashset"
	b "github.com/WiggidyW/weve-esi/client/configure/internal/bucket/writer"
)

type AuthHashSetWriterClient = b.SAC_BucketWriterClient[ahs.AuthHashSet]
