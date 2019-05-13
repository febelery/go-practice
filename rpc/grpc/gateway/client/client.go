package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "learn/rpc/grpc/gateway/proto"
	"log"
	"os"
	"strings"
)

func main() {
	var keyPath string
	baseDir, _ := os.Getwd()
	if strings.Contains(baseDir, "rpc") {
		keyPath = fmt.Sprintf("%s/../conf/certs/server.pem", baseDir)
	} else {
		keyPath = fmt.Sprintf("%s/rpc/grpc/gateway/conf/certs/server.pem", baseDir)
	}

	creds, err := credentials.NewClientTLSFromFile(keyPath, "grpc.abc")
	if err != nil {
		log.Printf("Failed to create TLS credentials %v", err)
	}

	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(creds))
	defer conn.Close()

	if err != nil {
		log.Println(err)
	}

	c := pb.NewHelloWorldClient(conn)
	context := context.Background()
	body := &pb.HelloWorldRequest{
		Referer: "Hello Grpc",
	}

	r, err := c.SayHelloWorld(context, body)
	if err != nil {
		log.Println(err)
	}

	log.Println(r.Message)
}
