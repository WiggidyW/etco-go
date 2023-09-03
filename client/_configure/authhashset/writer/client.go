package writer

import (
	b "github.com/WiggidyW/etco-go-bucket"

	bwriter "github.com/WiggidyW/etco-go/client/configure/internal/bucket/writer"
)

type AuthHashSetWriterClient = bwriter.SAC_BucketWriterClient[b.AuthHashSet]
