package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/kjmkznr/terraform-provider-elasticsearch"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: elasticsearch.Provider,
	})
}
