package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoConnection = ConnectToDB()
var clientOptions = options.Client().ApplyURI("mongodb+srv://user:pa$$w0rd@cluster0.rlmwzkd.mongodb.net/test")

func ConnectToDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	log.Printf("Successful connection to MongoDB")
	return client
}

func CheckConnection() bool {
	err := MongoConnection.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}
	return true
}
