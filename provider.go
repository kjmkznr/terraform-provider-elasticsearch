package elasticsearch

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	// ElasticsearchURL endpoint
	ElasticsearchURL = "ELASTICSEARCH_URL"
	// ElasticsearchVersion major version, it should be prefixed "v"
	ElasticsearchVersion = "ELASTICSEARCH_VERSION"
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
			"es_version": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(ElasticsearchVersion, "v5"),
				Description: "elasticsearch version",
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
		URL:     d.Get("url").(string),
		Version: d.Get("es_version").(string),
	}

	return config.NewClient()
}
