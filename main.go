package main

import (
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"github.com/ubozov/grpc-atlant/products"
	pb "github.com/ubozov/grpc-atlant/proto/products/v1"
)

func main() {

	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	listener, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatalln("Failed to create listener:", err)
	}

	var opts []grpc.ServerOption

	srv := grpc.NewServer(opts...)
	pb.RegisterProductServiceServer(srv, &products.Service{})
	reflection.Register(srv)

	log.Infoln("Serving gRPC on:", listener.Addr().String())
	err = srv.Serve(listener)
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
