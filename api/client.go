package main

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	api "github.com/olegsmetanin/golang-grpc-rest-gorm-example/api/proto"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"github.com/olegsmetanin/golang-grpc-rest-gorm-example/api/cert"
	"crypto/x509"
	"fmt"
)

const (
	port = 10000
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.

	var err error


	demoCertPool := x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM([]byte(cert.Cert))
	if !ok {
		panic("bad certs")
	}
	demoAddr := fmt.Sprintf("localhost:%d", port)


	var opts []grpc.DialOption


	creds := credentials.NewClientTLSFromCert(demoCertPool, demoAddr)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(demoAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()


	c := api.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &api.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}