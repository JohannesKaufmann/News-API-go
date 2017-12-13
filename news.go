package news

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

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
	myClient = &http.Client{Timeout: 10 * time.Second}
)

// getJSON is fetching json from an api endpoint.
func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
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

func fetch(url string, opt interface{}) (networkResult, *Exception) {
	var res networkResult

	// convert the options struct to a query parameter string
	v, err := query.Values(opt)
	if err != nil {
		return res, &Exception{
			Code:    "[assembling query string]",
			Message: err.Error(),
		}
	}

	// attach the query parameter to the url
	url = url + "?" + v.Encode()

	err = getJSON(url, &res)
	if err != nil {
		return res, &Exception{
			Code:    "[requesting json]",
			Message: err.Error(),
		}
	}
	if res.Status == "error" {
		return res, &Exception{
			Code:    res.Code,
			Message: res.Message,
		}
	}

	return res, nil
}
