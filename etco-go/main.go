package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/proto"
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

	// initialize the service, which implements all protobuf methods
	service := service.NewService()

	// create the GRPC server and register the service
	var grpcServer *grpc.Server
	if build.DEV_MODE {
		creds, err := credentials.NewServerTLSFromFile(
			"cert.pem",
			"key.pem",
		)
		if err != nil {
			panic(err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		grpcServer = grpc.NewServer()
	}
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
