package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"learn/rpc/grpc/features/proto"
	"log"
	"time"
)

var addr = flag.String("addr", ":50051", "the addr to connect to")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := featuresProto.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	if err := sendMessage(stream, "hello"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}
	if err := sendMessage(stream, "world"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}

	recvMessage(stream, codes.OK)
	recvMessage(stream, codes.OK)

	fmt.Println("cancelling context")
	cancel()

	sendMessage(stream, "closed")
	recvMessage(stream, codes.Canceled)

}

func recvMessage(stream featuresProto.Echo_BidirectionalStreamingEchoClient, wantErrCode codes.Code) {
	res, err := stream.Recv()
	if status.Code(err) != wantErrCode {
		log.Fatalf("stream.Recv() = %v, %v; want _, status.Code(err)=%v", res, err, wantErrCode)
	}
	if err != nil {
		fmt.Printf("stream.Recv() returned expected error %v\n", err)
		return
	}

	fmt.Printf("received message %q\n", res.GetMessage())
}

func sendMessage(stream featuresProto.Echo_BidirectionalStreamingEchoClient, msg string) error {
	fmt.Printf("sending message %q\n", msg)
	return stream.Send(&featuresProto.EchoRequest{Message: msg})
}
