package model

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"simple-web-golang/util"
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
	End      Date    `json:"begin"`
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

func ArtistsByName(ctx context.Context, name string) ([]Artist, error) {
	ret := make([]Artist, 0)
	cfg, err := util.ConfFrom(ctx)
	if err != nil {
		return ret, err
	}
	cli, err := util.MongoFrom(ctx)
	if err != nil {
		return ret, err
	}

	coll := cli.Database(cfg.Db.DbName).Collection("artists")
	result, err := coll.Find(ctx, bson.D{{"name", name}})
	if err != nil {
		return ret, err
	}

	for result.Next(ctx) {
		var a Artist
		if err := result.Decode(&a); err != nil {
			continue
		}
		ret = append(ret, a)
	}

	return ret, nil
}
