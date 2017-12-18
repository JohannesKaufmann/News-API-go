package news

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type ClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	if c.DoFunc != nil {
		return c.DoFunc(req)
	}

	return &http.Response{}, nil
}

func testCacheHeader(t *testing.T, expected bool, callback func() *Exception) {
	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			included := req.Header.Get("X-No-Cache") == "true"
			if included != expected {
				t.Fatal("without 'ForceFreshData' the header 'X-No-Cache' should not be added OR the other way around")
				// t.Fatal("with 'ForceFreshData' the header 'X-No-Cache' should not be missing")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
			}, nil
		},
	}

	err := callback()

	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "body is nil") {
		t.Fatal("expected the error to be about the missing body")
	}
}
func testAPIKey(t *testing.T, expected string, callback func() *Exception) {
	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			key := req.URL.Query().Get("apiKey")
			if key != expected {
				t.Fatal("expected a different api key: ", key, " != ", expected)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
			}, nil
		},
	}

	err := callback()

	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "body is nil") {
		t.Fatal("expected the error to be about the missing body")
	}
}
func testHeaders(t *testing.T, expected map[string]string, callback func() *Exception) {
	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			headers := make(map[string]string)
			for key := range req.Header {
				headers[key] = req.Header.Get(key)
			}

			if !reflect.DeepEqual(headers, expected) {
				t.Fatal("the headers on the request are not the same")
			}

			return &http.Response{
				StatusCode: http.StatusOK,
			}, nil
		},
	}

	err := callback()

	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "body is nil") {
		t.Fatal("expected the error to be about the missing body")
	}
}

// fmt.Printf(" %+v\n", req.Header)
// Body:       ioutil.NopCloser(bytes.NewBufferString("Hello World")),

func TestTopHeadlines_WithCache(t *testing.T) {
	testCacheHeader(t, false, func() *Exception {
		opt := TopHeadlinesOptions{}
		_, _, err := TopHeadlines(opt)
		return err
	})
}
func TestTopHeadlines_WithoutCache(t *testing.T) {
	testCacheHeader(t, true, func() *Exception {
		opt := TopHeadlinesOptions{
			ForceFreshData: true,
		}
		_, _, err := TopHeadlines(opt)
		return err
	})
}
func TestTopHeadlines_LocalAPIKey(t *testing.T) {
	testAPIKey(t, "abc", func() *Exception {
		opt := TopHeadlinesOptions{
			APIKey: "abc",
		}
		_, _, err := TopHeadlines(opt)
		return err
	})
}
func TestTopHeadlines_GlobalAPIKey(t *testing.T) {
	APIKey = "abc"

	testAPIKey(t, "abc", func() *Exception {
		opt := TopHeadlinesOptions{}
		_, _, err := TopHeadlines(opt)
		return err
	})
}
func TestTopHeadlines_Headers(t *testing.T) {
	Headers = map[string]string{
		"key":   "value",
		"key-2": "value-2",
	}

	// with uppercase keys
	h := map[string]string{
		"Key":   "value",
		"Key-2": "value-2",
	}
	testHeaders(t, h, func() *Exception {
		opt := TopHeadlinesOptions{}
		_, _, err := TopHeadlines(opt)
		return err
	})
}
func TestTopHeadlines_Articles(t *testing.T) {
	json := `
	{
		"status":"ok",
		"totalResults": 2,
		"articles":[
			{
				"source":{
					"id":"bbc-news",
					"name":"BBC News"
				},
				"author":"BBC News",
				"title":"Title 1",
				"description":"Description 1",
				"url":"http://www.bbc.co.uk/news/1",
				"urlToImage":"image url 1",
				"publishedAt":"2017-12-18T16:27:39Z"
			},
			{
				"source":{
					"id":"bbc-news",
					"name":"BBC News"
				},
				"author":"BBC News",
				"title":"Amtrak train plunges on to US highway",
				"description":"There are reports of fatalities after the train fell from a bridge on to the road in Washington state.",
				"url":"http://www.bbc.co.uk/news/2",
				"urlToImage":"image url 2",
				"publishedAt":null
			}
		]
	}
	`

	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(json)),
			}, nil
		},
	}

	opt := TopHeadlinesOptions{}
	articles, info, err := TopHeadlines(opt)

	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 2 {
		t.Fatal("expected 2 articles but got ", len(articles))
	}
	if info.TotalResults != 2 {
		t.Fail()
	}

}
