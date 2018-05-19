//go:generate protoc -I echo --go_out=plugins=grpc:echo echo/echo.proto

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	mwgrpc "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"

	pb "github.com/kublr/workshop-microservice-build-pipeline-aggregator/pkg/aggregator"
)

var (
	port        = flag.Int("port", 11000, "The server port")
	colorerAddr = flag.String("colorer", "127.0.0.1:10000", "Colorer address in the format of host:port")
)

func main() {
	// parse flags
	flag.Parse()

	// prepare requested port listener
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// prepare server options including interceptors
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(mwgrpc.ChainStreamServer(
			otgrpc.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(mwgrpc.ChainUnaryServer(
			otgrpc.UnaryServerInterceptor(),
		)),
	}

	// create gRPC server
	grpcServer := grpc.NewServer(opts...)

	// register gRPC procedures handler
	pb.RegisterAggregatorServer(grpcServer, pb.NewServer(*colorerAddr))

	// start server
	grpcServer.Serve(lis)
}
