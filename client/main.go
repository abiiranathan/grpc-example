package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"rpc/proto/calculator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadCredentials() (credentials.TransportCredentials, error) {
	ca := x509.NewCertPool()
	caPath := "certs/rootca.crt"
	caCartBytes, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	ca.AppendCertsFromPEM(caCartBytes)

	certs, err := tls.LoadX509KeyPair("certs/certfile.crt", "certs/keyfile.key")
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certs},
		RootCAs:      ca,
	}), nil
}

func main() {
	creds, err := loadCredentials()
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	// create client
	client := calculator.NewCalculatorClient(conn)

	// add
	addReq := &calculator.CalcRequest{
		A: 10,
		B: 20,
	}
	addRes, err := client.Add(context.Background(), addReq)
	if err != nil {
		log.Fatalf("failed to add: %v", err)
	}

	log.Printf("Add result: %v", addRes.Result)

	// subtract
	subtractReq := &calculator.CalcRequest{
		A: 20,
		B: 10,
	}

	subtractRes, err := client.Subtract(context.Background(), subtractReq)
	if err != nil {
		log.Fatalf("failed to subtract: %v", err)
	}

	log.Printf("Subtract result: %v", subtractRes.Result)

	// multiply
	multiplyReq := &calculator.CalcRequest{
		A: 10,
		B: 20,
	}

	multiplyRes, err := client.Multiply(context.Background(), multiplyReq)
	if err != nil {
		log.Fatalf("failed to multiply: %v", err)
	}
	log.Printf("Multiply result: %v", multiplyRes.Result)

	// divide
	divideReq := &calculator.CalcRequest{
		A: 20,
		B: 10,
	}

	divideRes, err := client.Divide(context.Background(), divideReq)
	if err != nil {
		log.Fatalf("failed to divide: %v", err)
	}

	log.Printf("Divide result: %v", divideRes.Result)

}
