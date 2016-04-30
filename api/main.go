package main

import (
	//"log"
	//"net"
	//"net/http"
	//"strings"
	"github.com/gengo/grpc-gateway/runtime"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	api "github.com/olegsmetanin/golang-grpc-rest-gorm-example/api/proto"
	"github.com/olegsmetanin/golang-grpc-rest-gorm-example/api/cert"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"net/http"
	"strings"
	"net"
	"log"
)

const (
	port = 10000
)

//
//// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
//// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

//
//// server is used to implement helloworld.GreeterServer.
type grpcService struct{}
//
//// SayHello implements helloworld.GreeterServer
func (s *grpcService) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	log.Print(in.Name)
	return &api.HelloReply{Message: "Hello " + in.Name}, nil
}
//
func newServer() *grpcService {
	return new(grpcService)
}

func main() {

	//fmt.Print(cert.Cert)

	var err error
	pair, err := tls.X509KeyPair([]byte(cert.Cert), []byte(cert.Key))
	if err != nil {
		panic(err)
	}
	demoKeyPair := &pair
	demoCertPool := x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM([]byte(cert.Cert))
	if !ok {
		panic("bad certs")
	}
	demoAddr := fmt.Sprintf("localhost:%d", port)



	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(demoCertPool, demoAddr))}

	grpcServer := grpc.NewServer(opts...)
	api.RegisterGreeterServer(grpcServer, newServer())
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: demoAddr,
		RootCAs:    demoCertPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux()
	err = api.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, demoAddr, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
	}

	mux.Handle("/", gwmux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    demoAddr,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*demoKeyPair},
		},
	}

	fmt.Printf("grpc on port: %d\n", port)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


}