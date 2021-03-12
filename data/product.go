package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Price        float64            `bson:"price,omitempty"`
	Counter      int32              `bson:"counter,omitempty"`
	LastModified time.Time          `bson:"lastModified,omitempty"`
}
