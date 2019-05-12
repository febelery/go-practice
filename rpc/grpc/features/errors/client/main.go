package main

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"learn/rpc/grpc/helloworld/proto"
	"log"
	"os"
	"time"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	c := helloworld.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: "world"})
	if err != nil {
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				log.Printf("quota failure: %s", info)
			default:
				log.Printf("unexpected type: %s", info)
			}
		}
		os.Exit(1)
	}

	log.Printf("Greeting: %s", r.Message)
}
