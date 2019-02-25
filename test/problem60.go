package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"github.com/go-redis/redis"
	"log"
	"os"
	"sync"
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

const (
	keyPrefix = "artist::"
)

func main() {
	fname := flag.String("fname", "artist.json.gz", "input file")
	flag.Parse()
	log.Println(*fname)

	//log.Print(artists[:10])
	kvs := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer kvs.Close()
	artists := ReadFile(*fname)

	log.Printf("complete to read: %d record", len(artists))

	dpc := NewDispatcher(20, kvs)
	for _, a := range artists {
		dpc.StoreArtist(a)
	}
	dpc.Wait()
	dpc.Close()
}

func ReadFile(fname string) []Artist {
	fp, err := os.Open(fname)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer fp.Close()
	gz, err := gzip.NewReader(fp)
	if err != nil {
		log.Fatalf("failed to open reader: %v", err)
	}
	defer gz.Close()

	artists := make([]Artist, 0, 100)
	scanner := bufio.NewScanner(gz)
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

type Dispatcher struct {
	wg        sync.WaitGroup
	NumWorker int
	queue     chan Artist
}

func NewDispatcher(n int, kvs *redis.Client) *Dispatcher {
	ret := &Dispatcher{
		wg:        sync.WaitGroup{},
		NumWorker: n,
		queue:     make(chan Artist, 1000),
	}

	for i := 0; i < n; i++ {
		go func(idx int) {
			task := 0
			defer log.Printf("%d worker done %d tasks", idx, task)
			for {
				select {
				case a, more := <-ret.queue:
					if !more {
						return
					}
					kvs.Set(keyPrefix+a.Name, a.Area, 0)
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

func (d *Dispatcher) StoreArtist(artist Artist) {
	d.wg.Add(1)
	d.queue <- artist
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) Close() {
	close(d.queue)
}
