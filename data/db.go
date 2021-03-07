package data

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
func (db DB) ProductCollection() *mongo.Collection {
	return db.client.Database(db.dbName).Collection("products")
}
