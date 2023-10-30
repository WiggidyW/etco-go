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
	ESI_USER_AGENT    string = os.Getenv("ESI_USER_AGENT")
	BUCKET_CREDS_JSON string = os.Getenv("BUCKET_CREDS_JSON")
	BUCKET_NAMESPACE  string = os.Getenv("BUCKET_NAMESPACE")
	SKIP_SDE          bool   = os.Getenv("SKIP_SDE") == "true"
	SKIP_CORE         bool   = os.Getenv("SKIP_CORE") == "true"
)

func main() {
	if BUCKET_CREDS_JSON == "" {
		log.Fatal("BUCKET_CREDS_JSON is empty")
	} else if ESI_USER_AGENT == "" {
		log.Print("ESI_USER_AGENT is empty")
	}
	updateResult, err := updater.Update(
		b.NewBucketClient(BUCKET_NAMESPACE, []byte(BUCKET_CREDS_JSON)),
		&http.Client{},
		ESI_USER_AGENT,
		SKIP_SDE,
		SKIP_CORE,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", updateResult)
}
