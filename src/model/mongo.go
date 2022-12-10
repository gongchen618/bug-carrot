package model

import (
	"bug-carrot/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"os"
)

var (
	mongoClient *mongo.Client
	dbName      string
)

func (m *model) Close() {
	// DO NOTHING
}

//connectMongo helps to connect this program to the mongo internet
func connectMongo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	MongoUser, ok1 := os.LookupEnv("MongoUser")
	MongoPwd, ok2 := os.LookupEnv("MongoPwd")
	var ok3 bool
	dbName, ok3 = os.LookupEnv("DBName")
	if !ok1 || !ok2 || !ok3 {
		log.Panic("mongo config required: set environment for MongoUser, MongoPwd, DBName")
		return
	}

	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s",
		MongoUser, MongoPwd, config.C.MongoDB.Host, dbName)
	fmt.Println(mongoUri)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Panic("mongo connection failed:", err.Error())
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic("mongo ping failed:", err.Error())
		return
	}

	mongoClient = client
}

func getMongoDataBase(ctx context.Context) *mongo.Database {
	err := mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Panic("mongo ping failed:", err)
	}

	return mongoClient.Database(dbName)
}

func init() {
	if !config.C.DatabaseUse {
		return
	}
	connectMongo()
}
