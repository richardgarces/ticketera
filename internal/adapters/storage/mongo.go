package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client *mongo.Client
}

func NewMongoRepository(uri string) *MongoRepository {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	return &MongoRepository{Client: client}
}
