package aggregator

import (
	"log"
	"time"

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
	// establish connection
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(colorerAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// create client stub
	colorerClient := colorerpb.NewColorerClient(conn)

	s := &aggregatorServer{
		colorerAddr:   colorerAddr,
		colorerConn:   conn,
		colorerClient: colorerClient,
	}

	return s
}
