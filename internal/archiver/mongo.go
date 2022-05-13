package archiver

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() (mongoClient *mongo.Client) {

	mongo_username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongo_password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	mongo_uri := "mongodb://" +
		mongo_username + ":" + mongo_password +
		"@127.0.0.1:27017/?authSource=admin"

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	failOnError(err, "Unable to create Mongo Client")

	mongoCtx, _ := context.WithTimeout(context.Background(), 1000*time.Second)
	err = mongoClient.Connect(mongoCtx)
	failOnError(err, "Unable to conect to Mongo")

	return mongoClient
}
