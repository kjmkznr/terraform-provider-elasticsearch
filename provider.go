package elasticsearch

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	ElasticsearchURL = "ELASTICSEARCH_URL"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(ElasticsearchURL, nil),
				Description: "url",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"elasticsearch_index_template": resourceElasticsearchIndexTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"elasticsearch_index_template_document": dataSourceElasticsearchIndexTemplateDocument(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		URL: d.Get("url").(string),
	}

	return config.NewClient()
}
