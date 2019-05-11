package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/status"
	"learn/rpc/grpc/features/proto"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct{}

func (*server) UnaryEcho(ctx context.Context, in *featuresProto.EchoRequest) (*featuresProto.EchoResponse, error) {
	fmt.Printf("UnaryEcho called with message %q\n", in.GetMessage())
	return &featuresProto.EchoResponse{Message: in.Message}, nil
}

func (*server) ServerStreamingEcho(*featuresProto.EchoRequest, featuresProto.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (*server) ClientStreamingEcho(featuresProto.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (*server) BidirectionalStreamingEcho(featuresProto.Echo_BidirectionalStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to server: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	featuresProto.RegisterEchoServer(s, &server{})

	s.Serve(lis)
}
