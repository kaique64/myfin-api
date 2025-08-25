package db

import (
	"context"
	"log"
	"time"

	"finix-api/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func Connect(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Não foi possível pingar o MongoDB:", err)
	}

	log.Println("Conectado ao MongoDB em", cfg.MongoURI)

	MongoClient = client
	MongoDatabase = client.Database(cfg.MongoDatabase)
}
