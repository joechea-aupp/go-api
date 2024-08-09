package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	return client, nil
}

func New(mongo *mongo.Client) {
	client = mongo
}

func collection(collection string) *mongo.Collection {
	return client.Database("go-api").Collection(collection)
}
