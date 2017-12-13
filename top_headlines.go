package news

// ArticleSource contains the identifier id and a display
// name for the source this article came from.
type ArticleSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Article contains all the information about an article
// like the author or title.
type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"`
}

// TopHeadlinesOptions contains the options that can be passed
// to the rest api. It gets converted to a query string and added
// to the url.
type TopHeadlinesOptions struct {
	// A comma-seperated string of identifiers (maximum 20) for the
	// news sources or blogs you want headlines from. Use the
	// '/sources' endpoint to locate these programmatically or look at the sources index.
	Sources []string `url:"sources"`

	// Keywords or phrase to search for.
	Query string `url:"q"`

	// The category you want to get headlines for.
	Category string `url:"category"`

	// The 2-letter ISO-639-1 code of the language you want to get headlines for.
	Language string `url:"language"`

	// The 2-letter ISO 3166-1 code of the country you want to get headlines for.
	Country string `url:"country"`

	APIKey string `url:"apiKey"`
}

// TopHeadlines provides up to 10 live top and breaking headlines for
// a single source, or multiple sources. You can also search for
// current top headlines with keywords and filters. Articles are
// sorted by the source alphabetically, and then by the position they
// appear on the source's page (top to bottom).
// 		This endpoint is great for retrieving headlines for display
// 		on news tickers or similar.
func TopHeadlines(opt TopHeadlinesOptions) ([]Article, *Exception) {
	// the base url
	url := "https://newsapi.org/v2/top-headlines"

	res, err := fetch(url, opt)

	return res.Articles, err
}
