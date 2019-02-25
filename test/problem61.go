package main

import (
	"flag"
	"github.com/go-redis/redis"
	"log"
	"strings"
)

const (
	keyPrefix = "artist::"
)

func main() {
	artists := flag.String("artists", "", "artists")
	flag.Parse()
	names := strings.Split(*artists, ",")
	log.Println(names)

	kvs := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer kvs.Close()

	result := map[string]string{}
	for _, name := range names {
		name = strings.TrimSpace(name)
		if cmd := kvs.Get(keyPrefix + name); cmd.Err() != nil {
			result[name] = ""
		} else {
			area := cmd.Val()
			log.Print(area)
			result[name] = area
		}
	}
	log.Print(result)
}
