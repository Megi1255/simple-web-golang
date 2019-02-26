package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"simple-web-golang/config"
	"time"
)

func NewMongoClient(cfg *config.StorageConfig) *mongo.Client {
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

func FromContext(c context.Context) *mongo.Client {
	val := c.Value(config.KeyStorage)
	return val.(*mongo.Client)
}
