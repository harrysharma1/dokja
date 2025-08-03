package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectToMongo() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB not reachable: %s", err)
	}
	fmt.Println("Connected to MongoDB")
}
