package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"learn/rpc/grpc/features/proto"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct{}

func (*server) UnaryEcho(context.Context, *featuresProto.EchoRequest) (*featuresProto.EchoResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (*server) ServerStreamingEcho(*featuresProto.EchoRequest, featuresProto.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (*server) ClientStreamingEcho(featuresProto.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (*server) BidirectionalStreamingEcho(stream featuresProto.Echo_BidirectionalStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			fmt.Printf("server: error receiving form stream: %v\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}

		fmt.Printf("echoing message %q\n", in.Message)
		stream.Send(&featuresProto.EchoResponse{Message: in.Message})
	}
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at port %v\n", lis.Addr())

	s := grpc.NewServer()
	featuresProto.RegisterEchoServer(s, &server{})

	s.Serve(lis)
}
