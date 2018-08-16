package elasticsearch

import (
	"github.com/hashicorp/errwrap"

	v5 "gopkg.in/olivere/elastic.v5"
	v6 "gopkg.in/olivere/elastic.v6"
)

// Config ...
type Config struct {
	URL     string
	Version string
}

type Clients struct {
	V5Client *v5.Client
	V6Client *v6.Client
	Version  int
}

// NewClient ...
func (c *Config) NewClient() (*Clients, error) {
	var clients *Clients
	version := 5
	var err error
	if c.Version == "v6" {
		cli, e := v6.NewClient(
			v6.SetURL(c.URL),
			v6.SetSniff(false),
		)
		if e != nil {
			err = errwrap.Wrapf("es v6 configuration error: {{err}}", e)
		}
		version = 6
		clients = &Clients{
			V5Client: nil,
			V6Client: cli,
			Version:  version,
		}
	} else {
		cli, e := v5.NewClient(
			v5.SetURL(c.URL),
			v5.SetSniff(false),
		)
		if e != nil {
			err = errwrap.Wrapf("es v5 configuration error: {{err}}", e)
		}
		clients = &Clients{
			V5Client: cli,
			V6Client: nil,
			Version:  version,
		}
	}
	return clients, err
}
