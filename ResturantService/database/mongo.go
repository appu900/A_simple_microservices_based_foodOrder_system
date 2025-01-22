package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	log.Println("Connected to database , databaseType:MongoDB!")
	Client = client
	return nil
}




func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("resturant_db").Collection(collectionName)
}
