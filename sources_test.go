package news

import (
	"strings"
	"testing"
)

func TestSources_WithoutAPIKey(t *testing.T) {
	opt := Options{}
	_, err := Sources(opt)

	if err == nil {
		t.Fail()
	}
	if err.Error() != "[missing api key]" {
		t.Fatal("error string is different than expected")
	}
}
func TestSources_WithWrongAPIKey(t *testing.T) {
	opt := Options{
		APIKey: "xxx",
	}
	_, err := Sources(opt)
	if err == nil {
		t.Fatal("the request is supposed to fail but it didn't")
	}
	if err.Code != "apiKeyInvalid" {
		t.Fatal("the error code is not 'apiKeyInvalid'")
	}

	if !strings.Contains(err.Message, "invalid or incorrect") {
		t.Fatal("the error message does not contain 'invalid or incorrect'")
	}
}
