package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"learn/rpc/grpc/helloworld/proto"
	"log"
	"net"
	"time"
)

var ports = []string{":10001", ":10002", ":10003"}

type server struct{}

func (*server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + req.Name}, nil
}

type slowServer struct{}

func (s *slowServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	time.Sleep(time.Duration(200) * time.Millisecond)
	return &helloworld.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)

	go s.Serve(lis)
	defer s.Stop()

	for i := 0; i < 3; i++ {
		lis, err := net.Listen("tcp", ports[i])
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()

		s := grpc.NewServer()
		if i == 2 {
			helloworld.RegisterGreeterServer(s, &slowServer{})
		} else {
			helloworld.RegisterGreeterServer(s, &server{})
		}

		go s.Serve(lis)
	}

	select {}
}
