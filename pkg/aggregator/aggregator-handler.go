package aggregator

import (
	"log"
	"time"

	mwgrpc "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	colorerpb "github.com/kublr/workshop-microservice-build-pipeline-aggregator/pkg/colorer"
)

type aggregatorServer struct {
	colorerAddr   string
	colorerConn   *grpc.ClientConn
	colorerClient colorerpb.ColorerClient
}

// GetEcho returns the feature at the given point.
func (s *aggregatorServer) Aggregate(ctx context.Context, msg *AggregateRequest) (*AggregateResponse, error) {
	log.Printf("Server aggregator called with message (%v)", msg)

	// define client call context with timeout
	colorerClientCtx, colorerClientCancel := context.WithTimeout(ctx, 2*time.Second)
	defer colorerClientCancel()

	// call client, collect results
	ranges := []*ColorRange{}
	number := int(msg.GetNumber())
	if number <= 0 {
		number = 9
	}

	for i := 0; i < number; i++ {
		resp, err := s.colorerClient.GetColor(colorerClientCtx, &colorerpb.GetColorRequest{})
		if err != nil {
			// represent client call error
			resp = &colorerpb.GetColorResponse{
				Cold: 0,
				Hot:  0,
			}
		}
		ranges = append(ranges, &ColorRange{
			Cold: resp.Cold,
			Hot:  resp.Hot,
		})
	}

	return &AggregateResponse{
		Ranges: ranges,
	}, nil
}

func NewServer(colorerAddr string) AggregatorServer {
	// specify dependency connection parameters
	opts := []grpc.DialOption{
		// non-TLS connection
		grpc.WithInsecure(),

		// open tracing integration
		grpc.WithUnaryInterceptor(mwgrpc.ChainUnaryClient(
			otgrpc.UnaryClientInterceptor(),
		)),
		grpc.WithStreamInterceptor(mwgrpc.ChainStreamClient(
			otgrpc.StreamClientInterceptor(),
		)),
	}

	// establish connection
	conn, err := grpc.Dial(colorerAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// create client stub
	colorerClient := colorerpb.NewColorerClient(conn)

	// return new initialized aggregator server
	s := &aggregatorServer{
		colorerAddr:   colorerAddr,
		colorerConn:   conn,
		colorerClient: colorerClient,
	}

	return s
}
