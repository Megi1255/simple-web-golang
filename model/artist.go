package model

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simple-web-golang/config"
	"simple-web-golang/storage"
)

type Alias struct {
	Name     string `json:"name"`
	SortName string `json:"sort_name"`
}

type Tag struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type Rate struct {
	VotesCount int     `json:"votes-count"`
	Value      float32 `json:"value"`
}

type LifeSpan struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
	Ended bool   `json:"ended"`
}

type Artist struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Gender   string   `json:"gender"`
	SortName string   `json:"sort-name"`
	Area     string   `json:"area"`
	Type     string   `json:"type"`
	Country  string   `json:"country"`
	Aliases  []Alias  `json:"aliases"`
	LifeSpan LifeSpan `json:"life-span"`
	Tags     []Tag    `json:"tags"`
	Rating   Rate     `json:"rating"`
}

func (a *Artist) Pretty() string {
	b, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return ""
	}
	return string(b)
}

func FindArtists(ctx context.Context, filter bson.D, opts ...*options.FindOptions) ([]Artist, error) {
	ret := make([]Artist, 0)
	cfg := config.FromContext(ctx)
	cli := storage.FromContext(ctx)

	coll := cli.Database(cfg.Db.DbName).Collection("artists")
	result, err := coll.Find(ctx, filter, opts...)
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

func ArtistsByName(ctx context.Context, name string, limit int64, sort string, order int) ([]Artist, error) {
	opt := options.Find().SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	return FindArtists(
		ctx,
		//bson.D{{"name", bson.D{{"$regex", "^"+name+"$"}, {"$options", "i"}}}},
		bson.D{{"name", name}},
		opt,
	)
}

func ArtistsByArea(ctx context.Context, area string, offset int64, limit int64, sort string, order int) ([]Artist, error) {
	opt := options.Find().SetSkip(offset).SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	return FindArtists(
		ctx,
		bson.D{{"area", area}},
		opt,
	)
}

func ArtistByTag(ctx context.Context, tag string, offset int64, limit int64, sort string, order int) ([]Artist, error) {
	opt := options.Find().SetSkip(offset).SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	return FindArtists(
		ctx,
		bson.D{
			{"tags", bson.D{
				{"$elemMatch", bson.D{
					{"value", tag},
				}},
			}},
		},
		opt,
	)
}
