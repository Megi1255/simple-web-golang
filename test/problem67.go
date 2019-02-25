package main

import (
	"context"
	"encoding/json"
	"flag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strings"
	"time"
)

type Alias struct {
	Name     string `json:"name"`
	SortName string `json:"sort_name"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Date  int `json:"date"`
}

type Tag struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}

type Rate struct {
	Count int `json:"count"`
	Value int `json:"value"`
}

type Artist struct {
	Id       int64   `json:"id"`
	Gid      string  `json:"gid"`
	Name     string  `json:"name"`
	SortName string  `json:"sort_name"`
	Area     string  `json:"area"`
	Aliases  []Alias `json:"aliases"`
	Begin    Date    `json:"begin"`
	Tags     []Tag   `json:"tags"`
	Rating   Rate    `json:"rating"`
}

func (a *Artist) Pretty() string {
	b, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return ""
	}
	return string(b)
}

func main() {
	alias := flag.String("alias", "", "aritst aliases to find")
	flag.Parse()

	client := NewMongoClient("mongodb://localhost:27017")
	coll := client.Database("test").Collection("artists")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := coll.Find(ctx, bson.D{{"aliases", *alias}})
	if err != nil {
		log.Fatalf("failed to find: %v", err)
	}
	defer result.Close(ctx)
	for result.Next(ctx) {
		var a Artist
		if err := result.Decode(&a); err != nil {
			log.Fatalf("failed to decode: %v", err)
		}
		log.Print(a.Pretty())
	}

}

func NewMongoClient(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
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
