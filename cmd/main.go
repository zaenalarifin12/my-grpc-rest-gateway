package main

import (
	"context"
	"flag"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	bank_gw "github.com/zaenalarifin12/grpc-course/protogen/gateway/go/proto/bank"
	hello_gw "github.com/zaenalarifin12/grpc-course/protogen/gateway/go/proto/hello"
	resl_gw "github.com/zaenalarifin12/grpc-course/protogen/gateway/go/proto/resiliency"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9000", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	var opts []grpc.DialOption

	cred, err := credentials.NewClientTLSFromFile("ssl/ca.crt", "")

	if err != nil {
		log.Fatalln("Can't create client credentials : ", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(cred))
	if err := bank_gw.RegisterBankServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	if err := hello_gw.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	if err := resl_gw.RegisterResiliencyServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
