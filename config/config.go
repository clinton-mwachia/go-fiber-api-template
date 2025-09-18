package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping error:", err)
	}

	DB = client.Database(dbName)
	log.Println("✅ Connected to MongoDB:", dbName)
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
