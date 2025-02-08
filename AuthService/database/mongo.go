package database

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() error {
	// Replace with your MongoDB connection string.
	clientOptions := options.Client().ApplyURI("mongodb+srv://codingPeers:220720100141%40Code@codepeers.mpjti.mongodb.net/?retryWrites=true&w=majority&appName=codePeers")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}
	log.Println("Connected to auth database , databaseType:MongoDB!")
	Client = client
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("auth_db").Collection(collectionName)
}





