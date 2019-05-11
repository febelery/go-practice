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
	"strings"
	"time"
)

var port = flag.Int("port", 50051, "port number")

type server struct {
	client featuresProto.EchoClient
	cc     *grpc.ClientConn
}

func (s *server) UnaryEcho(ctx context.Context, req *featuresProto.EchoRequest) (*featuresProto.EchoResponse, error) {
	message := req.Message
	if strings.HasPrefix(message, "[propagate me]") {
		time.Sleep(800 * time.Millisecond)
		message = strings.TrimPrefix(message, "[propagate me]")

		return s.client.UnaryEcho(ctx, &featuresProto.EchoRequest{Message: message})
	}

	if message == "delay" {
		time.Sleep(1500 * time.Millisecond)
	}

	return &featuresProto.EchoResponse{Message: req.Message}, nil
}

func (s *server) ServerStreamingEcho(*featuresProto.EchoRequest, featuresProto.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "RPC unimplemented")
}

func (s *server) ClientStreamingEcho(featuresProto.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "RPC unimplemented")
}

func (s *server) BidirectionalStreamingEcho(stream featuresProto.Echo_BidirectionalStreamingEchoServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return status.Error(codes.InvalidArgument, "request message not received")
		}
		if err != nil {
			return err
		}

		message := req.Message
		if strings.HasPrefix(message, "[propagate me]") {
			time.Sleep(800 * time.Millisecond)
			message = strings.TrimPrefix(message, "[propagate me]")

			res, err := s.client.UnaryEcho(stream.Context(), &featuresProto.EchoRequest{Message: message})
			if err != nil {
				return err
			}

			stream.Send(res)
		}

		if message == "delay" {
			time.Sleep(1500 * time.Millisecond)
		}
		stream.Send(&featuresProto.EchoResponse{Message: message})
	}
}

func (s *server) Close() {
	s.cc.Close()
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server started at: %v", lis.Addr())

	echoServer := newEchoServer()
	defer echoServer.Close()

	grpcServer := grpc.NewServer()
	featuresProto.RegisterEchoServer(grpcServer, echoServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newEchoServer() *server {
	cc, err := grpc.Dial(fmt.Sprintf(":%d", *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &server{client: featuresProto.NewEchoClient(cc), cc: cc}
}
