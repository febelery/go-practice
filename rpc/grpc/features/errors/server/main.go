package main

import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"learn/rpc/grpc/helloworld/proto"
	"log"
	"net"
	"sync"
)

type server struct {
	mux   sync.Mutex
	count map[string]int
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.count[in.Name]++
	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "request limit exceeded.")
		ds, err := st.WithDetails(&errdetails.QuotaFailure{
			Violations: []*errdetails.QuotaFailure_Violation{{
				Subject:     fmt.Sprintf("name:%s", in.Name),
				Description: "Limit one greeting per person",
			}},
		})
		if err != nil {
			return nil, st.Err()
		}

		return nil, ds.Err()
	}

	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to server on: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
