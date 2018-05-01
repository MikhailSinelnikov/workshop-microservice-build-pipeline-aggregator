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

	"google.golang.org/grpc"

	pb "github.com/kublr/workshop-microservice-build-pipeline-aggregator/pkg/aggregator"
)

var (
	port        = flag.Int("port", 11000, "The server port")
	colorerAddr = flag.String("colorer", "127.0.0.1:10000", "Colorer address in the format of host:port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAggregatorServer(grpcServer, pb.NewServer(*colorerAddr))
	grpcServer.Serve(lis)
}
