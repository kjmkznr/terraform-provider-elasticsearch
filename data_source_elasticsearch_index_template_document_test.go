package elasticsearch

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccElasticsearchIndexTemplateDocument(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchIndexTemplateDocumentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateValue(
						"data.elasticsearch_index_template_document.test",
						"json",
						testAccElasticsearchIndexTemplateDocumentExpectedJSON,
					),
				),
			},
		},
	})
}

func testAccCheckStateValue(id, name, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		v := rs.Primary.Attributes[name]
		if v != value {
			return fmt.Errorf(
				"Value for %s is %s, not %s", name, v, value)
		}

		return nil
	}
}

var testAccElasticsearchIndexTemplateDocumentConfig = `
data "elasticsearch_index_template_document" "test" {
	template = "access_log-*"

	mapping {
		type = "nginx"

		_source {
			enabled = false
		}

		property {
			name   = "timestamp"
			format = "date_time_no_millis"
			type   = "date"
		}

		property {
			name   = "request_uri"
			index  = "not_analyzed"
			type   = "string"
		}
	}

	mapping {
		type = "apache"

		property {
			name   = "timestamp"
			format = "date_time_no_millis"
			type   = "date"
		}

		property {
			name   = "uri"
			index  = "not_analyzed"
			type   = "string"
		}
	}

	setting {
		number_of_shards = 1
	}
}
`
var testAccElasticsearchIndexTemplateDocumentExpectedJSON = `{
  "settings": {
    "number_of_shards": 1
  },
  "mappings": {
    "apache": {
      "properties": {
        "timestamp": {
          "format": "date_time_no_millis",
          "type": "date"
        },
        "uri": {
          "index": "not_analyzed",
          "type": "string"
        }
      }
    },
    "nginx": {
      "properties": {
        "request_uri": {
          "index": "not_analyzed",
          "type": "string"
        },
        "timestamp": {
          "format": "date_time_no_millis",
          "type": "date"
        }
      },
      "_source": {
        "enabled": false
      }
    }
  },
  "template": "access_log-*"
}`
