package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/WiggidyW/etco-go/bucket"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/proto"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/service"
)

var PORT = os.Getenv("PORT")

const DEFAULT_ADDR string = ":8080"

func getAddr() string {
	if PORT == "" {
		return DEFAULT_ADDR
	} else {
		return fmt.Sprintf(":%s", PORT)
	}
}

func main() {
	timeStart := time.Now()

	// initialize basal clients, upon which service inner clients are built
	cCache := cache.NewSharedClientCache(
		build.CCACHE_MAX_BYTES,
	)
	sCache := cache.NewSharedServerCache(
		build.SCACHE_ADDRESS,
	)
	rBucketClient := bucket.NewBucketClient(
		build.BUCKET_NAMESPACE,
		[]byte(build.BUCKET_CREDS_JSON),
	)
	rRDBClient := rdb.NewRemoteDBClient(
		[]byte(build.REMOTEDB_CREDS_JSON),
		build.REMOTEDB_PROJECT_ID,
	)
	httpClient := &http.Client{} // TODO: timeouts are defined per-method

	// initialize the service, which implements all protobuf methods
	service := service.NewService(
		cCache,
		sCache,
		rBucketClient,
		rRDBClient,
		httpClient,
	)

	// DELETE THIS DO NOT COMMIT //
	creds, err := credentials.NewServerTLSFromFile(
		"cert.pem",
		"key.pem",
	)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// create the GRPC server and register the service
	// grpcServer := grpc.NewServer()
	proto.RegisterEveTradingCoServer(grpcServer, service)

	listener, err := net.Listen("tcp", getAddr())
	if err != nil {
		panic(err)
	}

	// log the time it took to start the server
	go func() {
		logger.Info(fmt.Sprintf(
			"Server started on %s in %s",
			getAddr(),
			time.Since(timeStart),
		))
	}()

	grpcServer.Serve(listener)
}
