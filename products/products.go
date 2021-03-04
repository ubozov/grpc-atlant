package products

import (
	"context"

	proto "github.com/ubozov/grpc-atlant/proto/products/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service ...
type Service struct{}

// Fetch ...
func (s *Service) Fetch(ctx context.Context, request *proto.FetchRequest) (*emptypb.Empty, error) {
	// TODO

	return &emptypb.Empty{}, nil
}

// List ...
func (s *Service) List(ctx context.Context, request *proto.ListRequest) (*proto.ListResponse, error) {
	// TODO

	return &proto.ListResponse{}, nil
}
