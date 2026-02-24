package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	Collection *mongo.Collection
}
func mongoose(collectionName string) *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI not set in .env")
	}

	clientOptions := options.Client().
		ApplyURI(uri)


	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	collection := client.
		Database("FooBar").
		Collection(collectionName)

	return &DB{
		Collection: collection,
	}
}