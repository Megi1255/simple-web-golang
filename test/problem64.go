package main

import (
	"bufio"
	_ "compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"os"
	"sync"
	"time"
)

/*
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
*/

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

type ArtistRaw struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	SortName string `json:"sort-name"`
	Area     struct {
		Name string `json:"name"`
	} `json:"area"`
	Type     string   `json:"type"`
	Country  string   `json:"country"`
	Aliases  []Alias  `json:"aliases"`
	LifeSpan LifeSpan `json:"life-span"`
	Tags     []Tag    `json:"tags"`
	Rating   Rate     `json:"rating"`
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

const (
	keyPrefix = "artist::"
)

func main() {
	fname := flag.String("fname", "artist.json", "input file")
	flag.Parse()
	log.Println(*fname)
	artists := ReadFile(*fname)
	log.Printf("complete to read: %v records", len(artists))

	client := NewMongoClient("mongodb://localhost:27017")
	coll := client.Database("new").Collection("artists")
	coll.Drop(context.Background())

	dpc := NewDispatcher(10, func(data interface{}) error {
		a := data.([]Artist)
		iSlice := make([]interface{}, len(a))
		for i, d := range a {
			iSlice[i] = d
		}
		col, err := coll.Clone()
		if err != nil {
			return err
		}
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		if _, err := col.InsertMany(ctx, iSlice); err != nil {
			return err
		}
		return nil
	})
	var i int
	for i = 0; i < len(artists)-100; i += 100 {
		dpc.Dispatch(artists[i : i+100])
	}
	dpc.Dispatch(artists[i:])
	dpc.Wait()
	dpc.Close()

	ims := []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"id", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
		/*
			{
				Keys:    bsonx.Doc{{"gid", bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true),
			},
		*/
		{
			Keys: bsonx.Doc{{"name", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"aliases.name", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"tags.value", bsonx.Int32(-1)}},
		},
		{
			Keys: bsonx.Doc{{"rating.value", bsonx.Int32(-1)}},
		},
	}
	if _, err := coll.Indexes().CreateMany(context.Background(), ims); err != nil {
		log.Fatalf("failed to create index: %v", err)
	}
}

func ReadFile(fname string) []Artist {
	fp, err := os.Open(fname)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer fp.Close()
	/*
		gz, err := gzip.NewReader(fp)
		if err != nil {
			log.Fatalf("failed to open reader: %v", err)
		}
		defer gz.Close()
	*/
	artists := make([]Artist, 0, 100)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		var artist Artist
		if err := json.Unmarshal([]byte(line), &artist); err != nil {
			log.Fatalf("failed to unmarshal: %v", err)
		}
		artists = append(artists, artist)
	}

	return artists
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

type Dispatcher struct {
	wg        sync.WaitGroup
	NumWorker int
	queue     chan interface{}
}

func NewDispatcher(n int, work func(interface{}) error) *Dispatcher {
	ret := &Dispatcher{
		wg:        sync.WaitGroup{},
		NumWorker: n,
		queue:     make(chan interface{}, 1000),
	}

	for i := 0; i < n; i++ {
		go func(idx int) {
			task := 0
			for {
				select {
				case a, more := <-ret.queue:
					if !more {
						log.Printf("%d worker done %d tasks", idx, task)
						return
					}
					if err := work(a); err != nil {
						log.Printf("%d worker error: %v", idx, err)
					}
					ret.wg.Done()
					if task += 1; task%1000 == 0 {
						log.Printf("%d worker: %d record", idx, task)
					}
				}
			}
		}(i)
	}
	return ret
}

func (d *Dispatcher) Dispatch(data interface{}) {
	d.wg.Add(1)
	d.queue <- data
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) Close() {
	close(d.queue)
}
