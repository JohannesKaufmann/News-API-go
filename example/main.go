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

func main() {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("error: the environment variable NEWS_API_KEY needs to be set.")
	}

	news.APIKey = apiKey
	news.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	news.Headers = map[string]string{
		"User-Agent": "Golang Client",
	}

	fmt.Println("\n - - sources - - ")
	sources()

	fmt.Println("\n - - top headlines - - ")
	topHeadlines()

	fmt.Println("\n - - everything - - ")
	everything()
}

func sources() {
	opt := news.SourcesOptions{
		ForceFreshData: true,
		Country:        "de",
	}

	sources, info, err := news.Sources(opt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info.Cached, info.Expires, info.Remaining, info.Date)
	fmt.Println("\nsources length: ", len(sources), "/", info.TotalResults)
	fmt.Printf("sources[0]: %+v\n", sources[0])
}

func topHeadlines() {
	opt := news.TopHeadlinesOptions{
		// APIKey: apiKey,
		Sources: []string{
			"bbc-news",
			"techcrunch",
		},
	}

	headlines, info, err := news.TopHeadlines(opt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nheadlines length: ", len(headlines), "/", info.TotalResults)
	fmt.Printf("headlines[0]: %+v\n", headlines[0])
}

func everything() {
	opt := news.EverythingOptions{
		// APIKey: apiKey,
		Query: "bitcoin",
	}

	everything, info, err := news.Everything(opt)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("\neverything length: ", len(everything), "/", info.TotalResults)
	fmt.Printf("everything[0]: %+v\n", everything[0])
}
