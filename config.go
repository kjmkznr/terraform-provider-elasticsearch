package elasticsearch

import elastic "gopkg.in/olivere/elastic.v5"

// Config ...
type Config struct {
	URL string
}

// NewClient ...
func (c *Config) NewClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(c.URL),
		elastic.SetSniff(false),
	)
	return client, err
}
