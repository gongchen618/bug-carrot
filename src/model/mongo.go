package model

import (
	"bug-carrot/src/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
)

var (
	mongoClient *mongo.Client
)

func (m *model) Close() {
	// DO NOTHING
}

//connectMongo helps to connect this program to the mongo internet
func connectMongo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s",
		config.C.MongoDB.Username, config.C.MongoDB.Password, config.C.MongoDB.Host, config.C.MongoDB.Database)
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

	return mongoClient.Database(config.C.MongoDB.Database)
}

func init() {
	if !config.C.DatabaseUse {
		return
	}
	connectMongo()
}
