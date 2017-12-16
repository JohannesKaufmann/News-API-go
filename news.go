package news

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// APIKey contains the api key required for every request
// to the news api. Gets added to the query string if its
// not already included.
var APIKey string

// Headers contains the request headers.
// for example "User-Agent": "Golang Client"
var Headers map[string]string

// Exception is a representation of a news exception.
type Exception struct {
	Code    string `json:"code"` // Error Code
	Message string `json:"message"`
}

// Error stringifies the error
func (e *Exception) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return e.Code + ": " + e.Message
}

var (
	// HTTPClient is the http client that is used for every
	// request. If you want to have a different timeout
	// overwrite the variable before using it.
	HTTPClient = &http.Client{Timeout: 10 * time.Second}
)

// getJSON is fetching json from an api endpoint.
func getJSON(url string, target interface{}, headers map[string]string) (http.Header, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// adding the headers that the user specified to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(target)

	return resp.Header, err
}

// networkResult is the standard result from the rest api.
type networkResult struct {
	Status string `json:"status"`

	Code    string `json:"code"`
	Message string `json:"message"`

	// - - either articles - - //
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`

	// - - or sources - - //
	Sources []Source `json:"sources"`
}

// ResponseInfo contains more information about the response
// like wether the request got cached or the nmuber of
// total articles.
type ResponseInfo struct {
	// INPUT
	ForceFreshData bool

	// wether the response is cached
	Cached bool

	// when the result will expire.
	Expires   string
	Remaining string

	// the date of the originally request that was then cached
	Date string

	TotalResults int
}

func fetch(url string, opt interface{}, forceFreshData bool) (networkResult, *ResponseInfo, *Exception) {
	var res networkResult

	// copy the values from the map over into the new one.
	// -> https://stackoverflow.com/a/23058707
	reqHeaders := make(map[string]string)
	for k, v := range Headers {
		reqHeaders[k] = v
	}

	// it the user specified that the cache should be circumvented
	// the header is added according to: https://newsapi.org/docs/caching
	if forceFreshData {
		reqHeaders["X-No-Cache"] = "true"
	}

	// convert the options struct to a query parameter string
	v, err := query.Values(opt)
	if err != nil {
		return res, nil, &Exception{
			Code:    "[assembling query string]",
			Message: err.Error(),
		}
	}

	// attach the query parameter to the url
	url = url + "?" + v.Encode()

	headers, err := getJSON(url, &res, reqHeaders)
	if err != nil {
		return res, nil, &Exception{
			Code:    "[requesting json]",
			Message: err.Error(),
		}
	}
	if res.Status != "ok" {
		return res, nil, &Exception{
			Code:    res.Code,
			Message: res.Message,
		}
	}

	isCached := headers.Get("X-Cached-Result") == "true"
	expires := headers.Get("X-Cache-Expires")
	remaining := headers.Get("X-Cache-Remaining")
	date := headers.Get("Date")

	if forceFreshData && isCached {
		log.Println("[news api] warning: you wanted fresh data but the api still returned cached data.")
	}
	if !forceFreshData && !isCached {
		log.Println("[DEBUG news api] info: you got fresh data although you did not want it.")
	}

	info := &ResponseInfo{
		ForceFreshData: forceFreshData,

		Cached:    isCached,
		Expires:   expires,
		Remaining: remaining,
		Date:      date,
	}
	fmt.Printf("%+v", info)

	return res, info, nil
}
