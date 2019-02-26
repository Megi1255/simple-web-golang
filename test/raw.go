package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func main() {
	ReadFile("artist/mbdump/artist", "artist.json")
}

func ReadFile(fname string, out string) {
	fp, err := os.Open(fname)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer fp.Close()

	outfp, err := os.Create(out)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer outfp.Close()

	cnt := 0
	reader := bufio.NewReader(fp)
	for {
		if cnt += 1; cnt%10000 == 0 {
			log.Printf("cnt: %d", cnt)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("failed to readline: %v", err)
			break
		}
		//log.Println(string(line))
		var raw ArtistRaw
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			log.Fatalf("failed to unmarshal: %v", err)
		}
		var a Artist
		a.Id, a.Name, a.Gender, a.SortName, a.Area, a.Type, a.Country, a.Aliases, a.LifeSpan, a.Tags, a.Rating =
			raw.Id, raw.Name, raw.Gender, raw.SortName, raw.Area.Name, raw.Type, raw.Country, raw.Aliases, raw.LifeSpan, raw.Tags, raw.Rating
		byteA, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("failed to marshal: %v", err)
		}
		if _, err := fmt.Fprintln(outfp, string(byteA)); err != nil {
			log.Fatalf("failed to write: %v", err)
		}
	}
}
