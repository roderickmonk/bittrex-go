package archiver

import (
	"context"
	// "fmt"
	"log"
	"os"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoCtx context.Context

// var err error

func NewMongoClient() {

	mongo_username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongo_password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	mongo_uri := "mongodb://" +
		mongo_username + ":" +
		mongo_password +
		"@127.0.0.1:27017/?authSource=admin"

	MongoClient, err = mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}
	MongoCtx, _ = context.WithTimeout(context.Background(), 1000*time.Second)
	err = MongoClient.Connect(MongoCtx)
	if err != nil {
		log.Fatal(err)
	}
}
