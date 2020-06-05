package main

import (
	"cloud/lib/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

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

var cfg = Config{
	Host:     os.Getenv("MONGO_DOCKER_HOST"),
	Port:     os.Getenv("MONGO_DOCKER_PORT"),
	User:     os.Getenv("MONGO_DOCKER_USER"),
	Password: os.Getenv("MONGO_DOCKER_PWD"),
	Database: os.Getenv("MONGO_DOCKER_DATABASE"),
}

func main() {
	// Connect()

	logger.Debug(cfg)
	db := GetDatabase(cfg)
	c := db.Collection("test")

	j := json.RawMessage(`{"a": "b"}`)
	res, _ := c.InsertOne(context.TODO(), bson.M{"name": "pi", "value": 3.14159, "json": j})
	id := res.InsertedID
	logger.Debug(id)
	cur, _ := c.Find(context.TODO(), bson.D{})

	for cur.Next(context.TODO()) {

		var md S
		err := cur.Decode(&md)
		if err != nil {
			logger.Error(err)
		}
		logger.Debug(md.String())

	}
}

func GetDatabase(cfg Config) *mongo.Database {

	uri := fmt.Sprintf(`mongodb://%s:%s@%s:%s`, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Panic(err)
	}

	// Check the connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Panic(err)
	}

	logger.Info("Connected to MongoDB!")

	return client.Database(cfg.Database)
}

type S struct {
	Name  string          `json:"name" bson:"name"`
	Value float64         `json:"value" bson:"value"`
	Json  json.RawMessage `json:"json" bson:"json"`
}

func (s S) String() string {
	by, _ := json.Marshal(s)
	return string(by)
}
