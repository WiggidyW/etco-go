package bucket

import (
	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
)

var (
	client *b.BucketClient
)

func init() {
	client = b.NewBucketClient(
		build.BUCKET_NAMESPACE,
		[]byte(build.BUCKET_CREDS_JSON),
	)
}
