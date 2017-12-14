package news

// Source contains a news publisher.
type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}

// SourcesOptions contains the options that can be passed
// to the rest api. It gets converted to a query string and added
// to the url.
type SourcesOptions struct {
	// if 'true' the header `X-No-Cache = true` will be added to
	// the request. Default: false (use cache)
	// -> https://newsapi.org/docs/caching
	ForceFreshData bool `url:"-"` // TODO: better name

	Category string `url:"category"`
	Language string `url:"language"`
	Country  string `url:"country"`
	APIKey   string `url:"apiKey"`
}

// Sources returns the subset of news publishers that top
// headlines (/v2/top-headlines) are available from. It's
// mainly a convenience endpoint that you can use to keep
// track of the publishers available on the API, and you
// can pipe it straight through to your users.
func Sources(opt SourcesOptions) ([]Source, *ResponseInfo, *Exception) {
	// the base url
	url := "https://newsapi.org/v2/sources"

	if opt.APIKey == "" && APIKey != "" {
		opt.APIKey = APIKey
	}

	res, info, err := fetch(url, opt, opt.ForceFreshData)

	// the response does not contain `res.TotalResults` so
	// I am setting it to the length of the array to
	// provided the same api to the user.
	if info != nil {
		info.TotalResults = len(res.Sources)
	}

	return res.Sources, info, err
}
