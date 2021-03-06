package products

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ubozov/grpc-atlant/data"
	proto "github.com/ubozov/grpc-atlant/proto/products/v1"
)

type product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Price        float64            `bson:"price,omitempty"`
	Counter      int32              `bson:"counter,omitempty"`
	LastModified time.Time          `bson:"lastModified,omitempty"`
}

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

// Start ...
func (s *Service) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption

	srv := grpc.NewServer(opts...)

	proto.RegisterProductServiceServer(srv, s)
	reflection.Register(srv)

	s.logger.Infoln("Serving gRPC on:", listener.Addr().String())

	return srv.Serve(listener)
}

// Fetch ...
func (s *Service) Fetch(ctx context.Context, request *proto.FetchRequest) (*emptypb.Empty, error) {
	//data, err := readCSVFromURL(request.Url)
	data, err := readCSVFromFile("./product.csv")
	if err != nil {
		return nil, err
	}

	if err != fetch(ctx, s.db, data) {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func fetch(ctx context.Context, db *data.DB, data [][]string) error {
	coll := db.ProductCollection()

	opts := options.Update().SetUpsert(true)

	for _, row := range data {
		price, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return err
		}

		if _, err := coll.UpdateOne(
			ctx,
			bson.D{
				{"name", row[0]},
			},
			bson.D{
				{"$set", bson.D{{"price", price}}}, //, {"lastModified", timestamppb.Now()}}},
				{"$currentDate", bson.D{{"lastModified", true}}},
				{"$inc", bson.M{"counter": 1}},
			},
			opts,
		); err != nil {
			return err
		}
	}
	return nil
}

// List ...
func (s *Service) List(request *proto.ListRequest, stream proto.ProductService_ListServer) error {
	coll := s.db.ProductCollection()

	ctx := context.Background()

	opts, filter, err := getSortAndPaginagionParams(request.PagingParam.Token, request.PagingParam.Limit,
		request.SortingParam.ColumnName, request.SortingParam.OrderType)
	if err != nil {
		return err
	}

	cursor, err := coll.Find(ctx, filter, &opts)
	if err != nil {
		return err
	}
	data := &product{}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		fmt.Println(cursor.ID(), cursor.Current)
		err := cursor.Decode(data)
		if err != nil {
			return err
		}

		stream.Send(&proto.ListResponse{
			Product: &proto.Product{
				Id:           data.ID.Hex(),
				Name:         data.Name,
				Price:        data.Price,
				Counter:      data.Counter,
				LastModified: timestamppb.New(data.LastModified),
			},
		})
	}

	return cursor.Err()
}

func getSortAndPaginagionParams(token string, limit int64, column string,
	order string) (options.FindOptions, interface{}, error) {
	opts := options.FindOptions{
		Limit: &limit,
	}
	if column == "" {
		column = "_id"
	}
	orderType := 1
	direction := "$gt"
	if order == "DESC" {
		orderType = -1
		direction = "$lt"
	}
	opts.SetSort(bson.D{{"_id", 1}, {column, orderType}})

	var filter interface{}
	var value interface{}

	if token != "" {
		s := strings.Split(token, "_")
		if len(s) != 3 {
			return opts, nil, fmt.Errorf("incorrect pagenation token")
		}

		switch s[0] {
		case "price":
			value, _ = strconv.ParseFloat(s[1], 64)
		case "counter":
			value, _ = strconv.ParseInt(s[1], 10, 64)
		}

		filter = bson.D{
			{s[0], bson.D{{direction, value}}},
			{"_id", bson.D{{"$ne", s[2]}}},
		}
	} else {
		filter = bson.D{}
	}

	return opts, filter, nil
}
