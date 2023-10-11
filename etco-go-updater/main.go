package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go-updater/updater"
)

var (
	ESI_USER_AGENT    = os.Getenv("ESI_USER_AGENT")
	BUCKET_CREDS_JSON = os.Getenv("BUCKET_CREDS_JSON")
)

func main() {
	if BUCKET_CREDS_JSON == "" {
		log.Fatal("BUCKET_CREDS_JSON is empty")
	} else if ESI_USER_AGENT == "" {
		log.Print("ESI_USER_AGENT is empty")
	}
	updateResult, err := updater.Update(
		b.NewBucketClient([]byte(BUCKET_CREDS_JSON)),
		&http.Client{},
		ESI_USER_AGENT,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", updateResult)
}
