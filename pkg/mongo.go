package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/karthik2304/xm-go-service/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (*mongo.Database, error) {
	// mongo
	clientOptions := options.Client().ApplyURI(configs.Settings.MONGO_ADDR)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("MongoDB Connection Error: %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("MongoDB Ping Failed: %v", err)
	}
	log.Println("Connected to MongoDB!")
	db := client.Database("companyDB")
	return db, nil
}
