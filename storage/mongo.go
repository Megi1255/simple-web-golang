package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func NewMongoClient(cfg *Config) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprint("mongodb://%s:%s", cfg.Host, cfg.Port)))
	if err != nil {
		log.Fatalf("failed to newclient: %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	return client
}
