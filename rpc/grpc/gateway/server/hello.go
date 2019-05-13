package server

import (
	"golang.org/x/net/context"
	pb "learn/rpc/grpc/gateway/proto"
)

type helloService struct{}

func (h helloService) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{
		Message: r.Referer,
	}, nil
}

func NewHelloService() *helloService {
	return &helloService{}
}
