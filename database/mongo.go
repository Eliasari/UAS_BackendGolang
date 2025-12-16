package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	if uri == "" || dbName == "" {
		log.Fatal("MONGO_URI atau MONGO_DB belum di-set di .env")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Gagal connect ke MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB tidak merespon:", err)
	}

	log.Println("MongoDB Connected ✅")

	return client.Database(dbName)
}
