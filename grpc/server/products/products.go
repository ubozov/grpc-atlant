package products

import (
	"context"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ubozov/grpc-atlant/data"
	proto "github.com/ubozov/grpc-atlant/proto/products/v1"
)

// Service ...
type Service struct {
	db     *data.DB
	logger grpclog.LoggerV2
}

// NewService creates a default instance.
func NewService(db *data.DB, logger grpclog.LoggerV2) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

// Fetch ...
func (s *Service) Fetch(ctx context.Context, request *proto.FetchRequest) (*emptypb.Empty, error) {
	data, err := readCSVFromURL(request.Url)
	//data, err := readCSVFromFile("./grpc/server/products/product.csv")
	if err != nil {
		return nil, err
	}

	if err != s.db.Fetch(ctx, data) {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// List ...
func (s *Service) List(request *proto.ListRequest, stream proto.ProductService_ListServer) error {
	ctx := context.Background()
	return s.db.List(ctx, request, stream)
}

// Log ...
func (s Service) Log(args ...interface{}) {
	s.logger.Infoln(args...)
}
