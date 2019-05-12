package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"learn/rpc/grpc/helloworld/proto"
	"log"
	"net"
	"os"
	"time"
)

const (
	defaultName = "world"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)

	go s.Serve(lis)
	defer s.Stop()

	r, rcleanup := manual.GenerateAndRegisterManualResolver()
	defer rcleanup()

	conn, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: ":10001"}, {Addr: ":10002"}, {Addr: ":10003"}}})

	c := helloworld.NewGreeterClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()

		r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: name})
		if err != nil {
			log.Printf("could not greet: %v", err)
		} else {
			log.Printf("Greeting: %s", r.Message)
		}
	}

	select {}
}
