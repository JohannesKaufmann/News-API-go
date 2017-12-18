package news

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestFetch_BadQueryString(t *testing.T) {
	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString("hello world")),
			}, nil
		},
	}

	url := "url"
	query := map[string]string{"should": "fail"}

	_, _, err := fetch(url, query, false)
	if err == nil {
		t.Fatal("expected error because of the bad query")
	}
	if !strings.Contains(err.Error(), "expects struct input") {
		t.Fatal(err)
	}
}
func TestFetch_NetworkError(t *testing.T) {
	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("some error")
		},
	}

	url := "url"
	query := struct {
		Name string `url:"name"`
	}{
		Name: "test",
	}

	_, _, err := fetch(url, query, false)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "some error") {
		t.Fatal(err)
	}
}

func TestFetch_StatusCode(t *testing.T) {
	HTTPClient = &ClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, nil
		},
	}

	url := "url"
	query := struct {
		Name string `url:"name"`
	}{
		Name: "test",
	}

	_, _, err := fetch(url, query, false)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "status code 400 != 200") {
		t.Fatal(err)
	}
}
