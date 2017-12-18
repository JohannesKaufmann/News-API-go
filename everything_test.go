package news

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestEverything_WithCache(t *testing.T) {
	testCacheHeader(t, false, func() *Exception {
		opt := EverythingOptions{}
		_, _, err := Everything(opt)
		return err
	})
}
func TestEverything_WithoutCache(t *testing.T) {
	testCacheHeader(t, true, func() *Exception {
		opt := EverythingOptions{
			ForceFreshData: true,
		}
		_, _, err := Everything(opt)
		return err
	})
}
func TestEverything_LocalAPIKey(t *testing.T) {
	testAPIKey(t, "abc", func() *Exception {
		opt := EverythingOptions{
			APIKey: "abc",
		}
		_, _, err := Everything(opt)
		return err
	})
}
func TestEverything_GlobalAPIKey(t *testing.T) {
	APIKey = "abc"

	testAPIKey(t, "abc", func() *Exception {
		opt := EverythingOptions{}
		_, _, err := Everything(opt)
		return err
	})
}
func TestEverything_Headers(t *testing.T) {
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
		opt := EverythingOptions{}
		_, _, err := Everything(opt)
		return err
	})
}

func TestEverything_Articles(t *testing.T) {
	json := `
	{
		"status":"ok",
		"totalResults":2,
		"articles":[
			{
				"source":{
					"id":null,
					"name":"Sueddeutsche.de"
				},
				"author":"Von Claus Hulverscheidt",
				"title":"Trump hilft Trump",
				"description":"Das große Bohei um die Steuerpläne verdeckt nicht, dass sie Reichen nutzen und Arbeitnehmern schaden.  Trump richtet seine Politik nicht zuletzt danach aus, ob sie dem Unternehmer und Privatmann Trump finanziell nutzt.",
				"url":"http://www.sueddeutsche.de/wirtschaft/us-steuerreform-trump-hilft-trump-1.3750403",
				"urlToImage":"http://polpix.sueddeutsche.com/staticassets/img/article/facebook-default.140428.jpg",
				"publishedAt":"2017-11-15T17:49:05Z"
			},
			{
				"source":{
					"id":"the-washington-post",
					"name":"The Washington Post"
				},
				"author":"(Greg Sargent)",
				"title":"Trump just rage-tweeted about a prominent African American again",
				"description":"Lavar ball, trump, donald trump, impeach donald trump, republican party, trump race baiting, trump attacks kneeling football players, gop",
				"url":"https://www.washingtonpost.com/blogs/plum-line/wp/2017/11/22/trump-just-rage-tweeted-about-a-prominent-african-american-again/",
				"urlToImage":"https://www.washingtonpost.com/rf/image_1484w/2010-2019/Wires/Online/2017-11-22/AP/Images/Trump_UCLA_Players_Basketball_76911.jpg-d1403.jpg?t=20170517",
				"publishedAt":"2017-11-22T15:22:15Z"
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

	opt := EverythingOptions{
		Query: "trump",
	}
	articles, info, err := Everything(opt)

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

func TestEverything_Error(t *testing.T) {
	json := `
	{
		"status":"error",
		"code":"parametersMissing",
		"message":"Required parameters are missing, the scope of your search is too broad. Please set any of the following required parameters and try again: q, sources, domains."
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

	opt := EverythingOptions{
		Query: "",
	}
	articles, _, err := Everything(opt)

	if err == nil {
		t.Fail()
	}
	if len(articles) != 0 {
		t.Fatal("expected an empty articles array")
	}
	// TODO: what should happen with info.TotalResults ?
}
