package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"rpc/proto/calculator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	calculator.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	return &calculator.CalcResponse{
		Result: req.A + req.B,
	}, nil
}

func (s *server) Divide(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	if req.B == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot divide by zero")
	}

	return &calculator.CalcResponse{
		Result: req.A / req.B,
	}, nil
}

func (s *server) Multiply(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	return &calculator.CalcResponse{
		Result: req.A * req.B,
	}, nil
}

func (s *server) Subtract(ctx context.Context, req *calculator.CalcRequest) (*calculator.CalcResponse, error) {
	return &calculator.CalcResponse{
		Result: req.A - req.B,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	tlscert, err := tls.LoadX509KeyPair("certs/certfile.crt", "certs/keyfile.key")
	if err != nil {
		log.Fatalf("failed to load key pair: %v", err)
	}

	s := server{}
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&tlscert)),
	)

	calculator.RegisterCalculatorServer(grpcServer, &s)

	// Enable reflection
	reflection.Register(grpcServer)

	defer grpcServer.GracefulStop()
	log.Println("gRPC Server is running on port https://localhost:9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
