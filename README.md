# News API SDK for Go

## Installation

`go get github.com/News-API-gh/News-API-go`

Make sure that you have an Api Key ready. I would recommend passing it as an environment variable. You can get a free key by [registering](https://newsapi.org/register).

## Features

* supports all [three endpoint](https://newsapi.org/docs/endpoints) with all their query parameters
* supports [disabling the cache](https://newsapi.org/docs/caching) to get fresh data
* supports adding custom headers to the request
* supports changing the http client

## Examples

### Getting all `Sources`

```golang
package main

import (
  news "github.com/News-API-gh/News-API-go"
)

func main() {
  opt := news.SourcesOptions{
    APIKey: "YOUR_API_KEY",
    Country: "de",
  }

  sources, info, err := news.Sources(opt)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(len(sources), "/", info.TotalResults)
  fmt.Printf("sources[0]: %+v\n", sources[0])
}
```

### Getting `Top Headlines`

```golang
package main

import (
  news "github.com/News-API-gh/News-API-go"
)

func main() {
  opt := news.TopHeadlinesOptions{
    APIKey: "YOUR_API_KEY",
    Sources: []string{
      "bbc-news",
      "techcrunch",
    },
  }

  headlines, info, err := news.TopHeadlines(opt)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(len(headlines), "/", info.TotalResults)
  fmt.Printf("headlines[0]: %+v\n", headlines[0])
}
```

### Getting `Everything`

```golang
package main

import (
  news "github.com/News-API-gh/News-API-go"
)

func main() {
  opt := news.EverythingOptions{
    APIKey: "YOUR_API_KEY",
    Query: "bitcoin",
  }

  everything, info, err := news.Everything(opt)
  if err != nil {
    log.Fatal(err.Error())
  }

  fmt.Println(len(everything), "/", info.TotalResults)
  fmt.Printf("everything[0]: %+v\n", everything[0])
}
```

### Setting the `Api Key` globally

You can also set the api key globally. This way you don't need to pass it to every function via the options parameter.

```golang
package main

import (
  news "github.com/News-API-gh/News-API-go"
)

func init() {
  apiKey := os.Getenv("NEWS_API_KEY")
  if apiKey == "" {
    log.Fatal("error: the environment variable NEWS_API_KEY needs to be set.")
  }

  // set the api key globally
  news.APIKey = apiKey
}
```

### Disabling the Cache

```golang
  opt := news.SourcesOptions{
    // if `ForceFreshData` is set to true the header
    // `X-No-Cache = true` is added to the request.
    ForceFreshData: true,
  }

  sources, info, err := news.Sources(opt)
  if err != nil {
    log.Fatal(err)
  }

  // The info struct always contains information about the caching
  // even if `ForceFreshData` is set to `false`.
  fmt.Println(info.Cached, info.Expires, info.Remaining, info.Date)

  fmt.Println(len(sources), "/", info.TotalResults)
  fmt.Printf("sources[0]: %+v\n", sources[0])
```

### Adding custom Headers and changing the Http Client

```golang
package main

import (
  news "github.com/News-API-gh/News-API-go"
)

func init() {
  news.HTTPClient = &http.Client{Timeout: 10 * time.Second}
  news.Headers = map[string]string{
    "User-Agent": "Golang Client",
  }
}
```

## TODO

* [ ] more tests

---

Coming soon... this is where our officially supported SDK for Go is going to live.

---

## Developers... we need you!

We need some help fleshing out this repo. If you're a Go dev with experience building client libraries or API wrappers, we're offering a reward of $250 to help us get started. For more info please email support@newsapi.org, or dive right in and send us a pull request.
