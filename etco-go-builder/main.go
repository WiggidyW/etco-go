package main

import (
	"context"
	"log"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go-builder/builder"
	"github.com/WiggidyW/etco-go-builder/builderenv"
)

func main() {
	if err := builderenv.ConvertAndValidate(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	bucketClient := b.NewBucketClient(
		[]byte(builderenv.BUCKET_CREDS_JSON),
	)

	err := builder.DownloadAndWrite(
		ctx,
		bucketClient,
		builderenv.GOB_FILE_DIR,
		builderenv.CONSTANTS_FILE_PATH,
	)
	if err != nil {
		log.Fatal(err)
	}
}
