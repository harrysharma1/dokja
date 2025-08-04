package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectToMongo() {
	var err error
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}
	log.Println("Connected to MongoDB")
}

func InsertWebNovel(novel WebNovel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := GetCollection().InsertOne(ctx, novel)
	return err
}

func GetCollection() *mongo.Collection {
	return client.Database("dokja").Collection("webnovels")
}

func FindAllWebNovels() ([]WebNovel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var novels []WebNovel

	cursor, err := GetCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var novel WebNovel
		if err := cursor.Decode(&novel); err != nil {
			log.Printf("Decode Error: %s", err)
			continue
		}
		novels = append(novels, novel)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return novels, nil
}

func FindWebNovelBasedOnUrlParam(urlPath string) (WebNovel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if urlPath != "" && urlPath[0] != '/' {
		urlPath = "/" + urlPath
	}

	var novel WebNovel

	err := GetCollection().FindOne(ctx, bson.M{"url_path": urlPath}).Decode(&novel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return novel, nil
		}
		return novel, err
	}
	return novel, nil
}
