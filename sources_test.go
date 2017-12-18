package news

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSources_WithCache(t *testing.T) {
	testCacheHeader(t, false, func() *Exception {
		opt := SourcesOptions{}
		_, _, err := Sources(opt)
		return err
	})
}
func TestSources_WithoutCache(t *testing.T) {
	testCacheHeader(t, true, func() *Exception {
		opt := SourcesOptions{
			ForceFreshData: true,
		}
		_, _, err := Sources(opt)
		return err
	})
}
func TestSources_LocalAPIKey(t *testing.T) {
	testAPIKey(t, "abc", func() *Exception {
		opt := SourcesOptions{
			APIKey: "abc",
		}
		_, _, err := Sources(opt)
		return err
	})
}
func TestSources_GlobalAPIKey(t *testing.T) {
	APIKey = "abc"

	testAPIKey(t, "abc", func() *Exception {
		opt := SourcesOptions{}
		_, _, err := Sources(opt)
		return err
	})
}
func TestSources_Headers(t *testing.T) {
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
		opt := SourcesOptions{}
		_, _, err := Sources(opt)
		return err
	})
}

func TestSources_Sources(t *testing.T) {
	json := `
	{
		"status":"ok",
		"sources":[
			{
				"id":"al-jazeera-english",
				"name":"Al Jazeera English",
				"description":"News, analysis from the Middle East and worldwide, multimedia and interactives, opinions, documentaries, podcasts, long reads and broadcast schedule.",
				"url":"http://www.aljazeera.com",
				"category":"general",
				"language":"en",
				"country":"us"
			},
			{
				"id":"bbc-news",
				"name":"BBC News",
				"description":"Use BBC News for up-to-the-minute news, breaking news, video, audio and feature stories. BBC News provides trusted World and UK news as well as local and regional perspectives. Also entertainment, business, science, technology and health news.",
				"url":"http://www.bbc.co.uk/news",
				"category":"general",
				"language":"en",
				"country":"gb"
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

	opt := SourcesOptions{}
	sources, info, err := Sources(opt)

	if err != nil {
		t.Fatal(err)
	}
	if len(sources) != 2 {
		t.Fatal("expected 2 sources but got ", len(sources))
	}
	if info.TotalResults != 2 {
		t.Fail()
	}

}
