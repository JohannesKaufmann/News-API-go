package news

// API stores basic information like the baseurl
// or apikey.
type API struct {
	BaseURL string
	APIKey  string
}

// NewClient creates a new News client.
func NewClient(apiKey string) *API {

	return &API{
		BaseURL: "https://newsapi.org/v2",
		APIKey:  apiKey,
	}
}
func (a *API) Sources(opt Options) ([]string, error) {

	return nil, nil
}
