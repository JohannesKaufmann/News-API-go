package news

// EverythingOptions contains the options that can be passed
// to the rest api. It gets converted to a query string and added
// to the url.
type EverythingOptions struct {
	// if 'true' the header `X-No-Cache = true` will be added to
	// the request. Default: false (use cache)
	// -> https://newsapi.org/docs/caching
	ForceFreshData bool `url:"-"`

	Query string `url:"q"`

	Sources []string `url:"sources"`

	Domains []string `url:"domains"`
	From    string   `url:"from"` // TODO: time.Time
	To      string   `url:"to"`

	Language string `url:"language"`
	SortBy   string `url:"sortBy"`
	Page     int    `url:"page,omitempty"`

	APIKey string `url:"apiKey"`
}

// Everything searches through millions of articles from over
// 5,000 large and small news sources and blogs. This includes
// breaking news as well as lesser articles.
//
// This endpoint suits article discovery and analysis, but can be
// used to retrieve articles for display, too.
func Everything(opt EverythingOptions) ([]Article, *ResponseInfo, *Exception) {
	// the base url
	url := "https://newsapi.org/v2/everything"

	if opt.APIKey == "" && APIKey != "" {
		opt.APIKey = APIKey
	}

	res, info, err := fetch(url, opt, opt.ForceFreshData)
	if info != nil {
		info.TotalResults = res.TotalResults
	}

	return res.Articles, info, err
}
