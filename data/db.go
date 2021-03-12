package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/ubozov/grpc-atlant/proto/products/v1"
)

// Config has a database connection configuration.
type Config struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

// DB has a configuration for database connect.
type DB struct {
	client *mongo.Client
	dbName string
}

// NewDB creates a new data connection handle.
func NewDB(ctx context.Context, conf Config) (*DB, error) {
	conn := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err == nil {
		err = client.Ping(ctx, readpref.Primary())
	}
	return &DB{client: client, dbName: conf.DBName}, err
}

// Close closes database connection.
func (db DB) Close(ctx context.Context) {
	db.client.Disconnect(ctx)
}

// ProductCollection returns products collection from the database.
func (db DB) productCollection() *mongo.Collection {
	return db.client.Database(db.dbName).Collection("products")
}

// Fetch ...
func (db DB) Fetch(ctx context.Context, data [][]string) error {
	coll := db.productCollection()

	opts := options.Update().SetUpsert(true)

	for _, row := range data {
		price, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return err
		}

		filter := bson.D{{"name", row[0]}}

		var p product
		err = coll.FindOne(ctx, filter).Decode(&p)

		op, inc := "$inc", 1
		if err != nil {
			if err != mongo.ErrNoDocuments {
				return err
			}
			op, inc = "$set", 0
		} else if p.Price == price {
			continue
		}

		if _, err := coll.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", bson.D{{"price", price}}},
				{"$currentDate", bson.D{{"lastModified", true}}},
				{op, bson.M{"counter": inc}},
			},
			opts,
		); err != nil {
			return err
		}
	}
	return nil
}

// List ...
func (db DB) List(ctx context.Context, req *proto.ListRequest, stream proto.ProductService_ListServer) error {
	coll := db.productCollection()

	opts, filter, err := getSortAndPaginagionOptions(req.PagingParam.Token, req.PagingParam.Limit,
		req.SortingParam.ColumnName, req.SortingParam.OrderType)
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

func getSortAndPaginagionOptions(token string, limit int64, column string,
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
			return opts, nil, fmt.Errorf("incorrect pagination token")
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
