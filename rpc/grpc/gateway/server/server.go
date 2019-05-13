package server

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"learn/rpc/grpc/gateway/pkg/util"
	pb "learn/rpc/grpc/gateway/proto"
	"net"
	"net/http"
)

var (
	Port        string
	SererName   string
	CertPemPath string
	CertKeyPath string
	EndPoint    string
	SwaggerDir  string

	tlsConfig *tls.Config
)

func Run() (err error) {
	EndPoint = ":" + Port
	tlsConfig = util.GetTLSConfig(CertPemPath, CertKeyPath)

	conn, err := net.Listen("tcp", EndPoint)
	if err != nil {
		log.Printf("TCP Listen err: %v\n", err)
	}

	srv := newServer(conn)

	log.Printf("gRPC and https listen on: %s\n", Port)

	if err = srv.Serve(util.NewTLSListener(conn, tlsConfig)); err != nil {
		log.Printf("ListenAndServer: %v\n", err)
	}

	return err
}

func newServer(conn net.Listener) *http.Server {
	grpcServer := newGrpc()
	gwmux, err := newGateway()

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	return &http.Server{
		Addr:      EndPoint,
		Handler:   util.GrpcHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}
}

func newGrpc() *grpc.Server {
	creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(creds),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterHelloWorldServer(server, NewHelloService())

	return server
}

func newGateway() (http.Handler, error) {
	ctx := context.Background()
	dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, SererName)
	if err != nil {
		return nil, err
	}

	dopts := []grpc.DialOption{
		grpc.WithTransportCredentials(dcreds),
	}

	gwmux := runtime.NewServeMux()
	if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
		return nil, err
	}

	return gwmux, nil
}
