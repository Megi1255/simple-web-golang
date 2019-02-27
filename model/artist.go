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

func FindArtists(ctx context.Context, filter bson.D, opts ...*options.FindOptions) (int64, []Artist, error) {
	ret := make([]Artist, 0)
	cfg := config.FromContext(ctx)
	cli := storage.FromContext(ctx)

	coll := cli.Database(cfg.Db.DbName).Collection("artists")
	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, ret, err
	}

	result, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return 0, ret, err
	}
	for result.Next(ctx) {
		var a Artist
		if err := result.Decode(&a); err != nil {
			continue
		}
		ret = append(ret, a)
	}

	return total, ret, nil
}

func ArtistsByName(ctx context.Context, name string, limit int64, sort string, order int) (int64, []Artist, error) {
	opt := options.Find().SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	/*
		filter := bson.D{{"name", bson.D{
			{"$regex", "^"+name}, {"$options", "i"},
			}},
		}
	*/
	filter := bson.D{{"name", name}}

	return FindArtists(ctx, filter, opt)
}

func ArtistsByArea(ctx context.Context, area string, offset int64, limit int64, sort string, order int) (int64, []Artist, error) {
	opt := options.Find().SetSkip(offset).SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	filter := bson.D{{"area", area}}

	return FindArtists(ctx, filter, opt)
}

func ArtistsByAlias(ctx context.Context, alias string, offset int64, limit int64, sort string, order int) (int64, []Artist, error) {
	opt := options.Find().SetSkip(offset).SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	filter := bson.D{
		{"aliases", bson.D{
			{"$elemMatch", bson.D{
				{"name", alias},
			}},
		}},
	}

	return FindArtists(ctx, filter, opt)
}

func ArtistByTag(ctx context.Context, tag string, offset int64, limit int64, sort string, order int) (int64, []Artist, error) {
	opt := options.Find().SetSkip(offset).SetLimit(limit)
	if sort != "" {
		opt.SetSort(bson.D{{sort, order}})
	}
	filter := bson.D{
		{"tags", bson.D{
			{"$elemMatch", bson.D{
				{"name", tag},
			}},
		}},
	}

	return FindArtists(ctx, filter, opt)
}
