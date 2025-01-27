package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var OrderClient *mongo.Client
var AuthClient *mongo.Client

func Connect() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderDBOptions := options.Client().ApplyURI("mongodb://mongodb-order:27017")
	orderClient, err := mongo.Connect(ctx, orderDBOptions)
	if err != nil {
		return err
	}
	err = orderClient.Ping(ctx, nil)
	if err != nil {
		return err
	}
	log.Println("Connected to order database , databaseType:MongoDB!")
	OrderClient = orderClient

	authDBOptions := options.Client().ApplyURI("mongodb://mongodb-primary:27017,mongodb-replica:27017/?replicaSet=rs0").SetReadPreference(readpref.Secondary())
	authClinet, err := mongo.Connect(ctx, authDBOptions)
	if err != nil {
		return err
	}
	err = authClinet.Ping(ctx, nil)
	if err != nil {
		return err
	}
	log.Println("connected to auth database , databaseType:MongoDB!")
	AuthClient = authClinet
	return nil
}

func Diconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if OrderClient != nil {
		if err := OrderClient.Disconnect(ctx); err != nil {
			return err
		}
	}

	if AuthClient != nil {
		if err := AuthClient.Disconnect(ctx); err != nil {
			return err
		}
	}
	return nil
}

func GetOrderCollection(collectionName string) *mongo.Collection {
	return OrderClient.Database("order_db").Collection(collectionName)
}

func GetAuthCollection(collectionName string) *mongo.Collection {
	return AuthClient.Database("auth_db").Collection(collectionName)
}
