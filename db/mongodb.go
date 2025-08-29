package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func CreateUniqueIndexOnChapters(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "webnovel_url_path", Value: 1},
			{Key: "chapter_number", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, error := collection.Indexes().CreateOne(ctx, indexModel)
	return error
}

func ConnectToMongo() {
	var err error
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	log.Println("Connected to MongoDB")
	errMakingUniqueChapterKey := CreateUniqueIndexOnChapters(GetCollectionChapters())
	if errMakingUniqueChapterKey != nil {
		log.Fatal("Failed to create unique index:", errMakingUniqueChapterKey)
	}
}

func InsertWebNovel(novel WebNovel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := GetCollectionNovels().InsertOne(ctx, novel)
	return err
}

func InsertChapter(chapter Chapter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, chapters, err := FindWebNovelBasedOnUrlParam(chapter.WebNovelUrlPath)
	if err != nil {
		return fmt.Errorf("failed to find chapters: %w", err)
	}
	for _, c := range chapters {
		if chapter.Number == c.Number {
			return fmt.Errorf("fhapter %d already exists", chapter.Number)
		}
	}
	_, errInserting := GetCollectionChapters().InsertOne(ctx, chapter)
	if errInserting != nil {
		return fmt.Errorf("failed to insert chapter: %w", errInserting)
	}
	return nil
}

func GetCollectionNovels() *mongo.Collection {
	return client.Database("dokja").Collection("webnovels")
}

func GetCollectionChapters() *mongo.Collection {
	return client.Database("dokja").Collection("chapters")
}

func FindAllWebNovels() ([]WebNovel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var novels []WebNovel

	cursor, err := GetCollectionNovels().Find(ctx, bson.M{})
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

func FindWebNovelBasedOnUrlParam(urlPath string) (WebNovel, []Chapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if urlPath != "" && urlPath[0] != '/' {
		urlPath = "/" + urlPath
	}

	var novel WebNovel
	var chapters []Chapter

	err := GetCollectionNovels().FindOne(ctx, bson.M{"url_path": urlPath}).Decode(&novel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return novel, chapters, nil
		}
		return novel, chapters, err
	}
	cursor, err := GetCollectionChapters().Find(ctx, bson.M{
		"webnovel_url_path": "/novels" + urlPath,
	})
	if err != nil {
		return novel, chapters, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var chapter Chapter
		if err := cursor.Decode(&chapter); err != nil {
			log.Printf("Decode Error: %s", err)
			continue
		}
		chapters = append(chapters, chapter)
	}

	if err := cursor.Err(); err != nil {
		return novel, chapters, err
	}
	return novel, chapters, nil
}

func FindChapterBasedOnUrlParam(urlPath string) (Chapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var chapter Chapter
	err := GetCollectionChapters().FindOne(ctx, bson.M{
		"url_path": urlPath,
	}).Decode(&chapter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return chapter, err
		}
		return chapter, err
	}

	return chapter, nil
}

func UpdateChapter(urlPath string, c Chapter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{"$set", bson.D{{"chapter_text", c.Text}}},
		{"$set", bson.D{{"chapter_number", c.Number}}},
		{"$set", bson.D{{"chapter_title", c.Title}}},
	}

	_, err := GetCollectionChapters().UpdateOne(ctx, bson.M{
		"url_path": urlPath,
	}, update)
	return err
}

func DeleteChapter(urlPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetCollectionChapters().DeleteOne(ctx, bson.M{
		"url_path": urlPath,
	})
	return err
}
