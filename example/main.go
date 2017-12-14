package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	news "github.com/JohannesKaufmann/News-API-go"
)

// TODO: the import needs to be rewritten

var apiKey string

func main() {
	apiKey = os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("error: the environment variable NEWS_API_KEY needs to be set.")
	}

	news.APIKey = apiKey

	sources()
	topHeadlines()
	everything()
}
func everything() {
	opt := news.EverythingOptions{
		// APIKey: apiKey,
		Query: "bitcoin",
	}

	everything, total, err := news.Everything(opt)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("\neverything length: ", len(everything), "/", total)
	fmt.Printf("everything[0]: %+v\n", everything[0])
}
func topHeadlines() {
	opt := news.TopHeadlinesOptions{
		// APIKey: apiKey,
		Sources: []string{
			"bbc-news",
			"techcrunch",
		},
	}

	headlines, total, err := news.TopHeadlines(opt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nheadlines length: ", len(headlines), "/", total)
	fmt.Printf("headlines[0]: %+v\n", headlines[0])
}
func sources() {
	news.HTTPClient = &http.Client{Timeout: 1 * time.Second}

	opt := news.SourcesOptions{
	// APIKey: apiKey,
	}

	sources, total, err := news.Sources(opt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nsources length: ", len(sources), "/", total)
	fmt.Printf("sources[0]: %+v\n", sources[0])
}
