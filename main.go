package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/service"
	"github.com/WiggidyW/etco-go/staticdb"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
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
	// initialize the logger
	go logger.InitLoggerCrashOnError()

	// initialize staticdb by loading .gob files, and crash on error
	go staticdb.LoadAllCrashOnError()

	// initialize the service, which implements all protobuf methods
	service := service.NewService()

	// create the GRPC server and register the service
	grpcServer := grpc.NewServer()
	proto.RegisterEveTradingCoServer(grpcServer, service)

	// wrap the server with GrpcWeb, enabling HTTP1.1 + Cors support
	// (HTTP2 still works - the wrapper just forwards non web requests)
	grpcWebServer := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(_ string) bool { return true }), // allow all origins
	)

	// create an HTTP server and serve the GRPCWeb server
	httpServer := &http.Server{
		Addr:    getAddr(), // 0.0.0.0:8080
		Handler: grpcWebServer,
	}
	httpServer.ListenAndServe()
}
