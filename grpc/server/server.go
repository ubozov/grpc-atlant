package server

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"github.com/ubozov/grpc-atlant/data"
	"github.com/ubozov/grpc-atlant/grpc/server/products"
	proto "github.com/ubozov/grpc-atlant/proto/products/v1"
)

// Start ...
func Start(db *data.DB, log grpclog.LoggerV2, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption

	srv := grpc.NewServer(opts...)

	s := products.NewService(db, log)

	proto.RegisterProductServiceServer(srv, s)
	reflection.Register(srv)

	s.Log("Serving gRPC on:", listener.Addr().String())

	return srv.Serve(listener)
}
