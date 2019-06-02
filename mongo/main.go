package main

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() {
	port := os.Getenv("MONGO_DOCKER_PORT")
	user := os.Getenv("MONGO_DOCKER_USER")
	pwd := os.Getenv("MONGO_DOCKER_PWD")
	logger.Debug(port, user, pwd)

	mongoURI := fmt.Sprintf("mongodb://localhost:%s/", port)

	opts := &options.ClientOptions{}
	opts.SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    "admin",
		Username:      user,
		Password:      pwd,
	})

	logger.Debug(opts.Auth)

	client, err := mongo.NewClient(opts.ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error(err)
		// return
	}

	DB := client.Database("admin")
	c := DB.Collection("test")
	res, err := c.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	id := res.InsertedID
	logger.Debug(id)
}

func main() {
	Connect()
}
