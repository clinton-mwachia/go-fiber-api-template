package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
	JWTTTLMin int
}

var (
	DB     *mongo.Database
	Client *mongo.Client
	Cfg    *Config
)

// a function to load env varibales
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found; falling back to environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "myapi_db"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	ttl := 60
	if v := os.Getenv("JWT_TTL_MIN"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			ttl = i
		}
	}

	Cfg = &Config{
		Port:      port,
		MongoURI:  mongoURI,
		MongoDB:   mongoDB,
		JWTSecret: jwtSecret,
		JWTTTLMin: ttl,
	}
}

// connect to the database
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

// DisconnectDB gracefully closes the MongoDB connection.
func DisconnectDB() {
	if Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Println("⚠️ Error disconnecting MongoDB:", err)
	} else {
		log.Println("👋 Disconnected from MongoDB")
	}
}
