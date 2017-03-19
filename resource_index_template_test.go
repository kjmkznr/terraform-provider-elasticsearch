package elasticsearch

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	elastic "gopkg.in/olivere/elastic.v5"
)

func TestAccIndexTemplate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndexTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIndexTemplateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foobar", "name", "template1"),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foobar", "template",
						regexp.MustCompile(`"template":"access_log-\*"`)),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foobar", "template",
						regexp.MustCompile(`"mappings":{"nginx":{"properties":{"timestamp":{`)),
				),
			},
		},
	})
}

func TestAccIndexTemplate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndexTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIndexTemplateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foobar", "name", "template1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIndexTemplateConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foobar", "name", "template1"),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foobar", "template",
						regexp.MustCompile(`"template":"access_log-\*"`)),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foobar", "template",
						regexp.MustCompile(`"timestamp":{`)),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foobar", "template",
						regexp.MustCompile(`"request_uri":{`)),
				),
			},
		},
	})
}

func testAccCheckIndexTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*elastic.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticsearch_index_template" {
			continue
		}

		svc := elastic.NewIndicesGetTemplateService(client)
		svc.Name(rs.Primary.ID)
		_, err := svc.Do(context.Background())
		if err == nil {
			return fmt.Errorf("Index template still exists")
		}
	}

	return nil
}

const testAccCheckIndexTemplateConfig_basic = `
resource "elasticsearch_index_template" "foobar" {
    name         = "template1"
	template = <<EOF
{
  "template": "access_log-*",
  "mappings": {
    "nginx": {
      "properties": {
        "timestamp": {
          "type": "date",
          "format":"date_time_no_millis"
        }
      }
    }
  }
}
EOF
}`

const testAccCheckIndexTemplateConfig_update = `
resource "elasticsearch_index_template" "foobar" {
    name         = "template1"
	template = <<EOF
{
  "template": "access_log-*",
  "mappings": {
    "nginx": {
      "properties": {
        "timestamp": {
          "type": "date",
          "format":"date_time_no_millis"
        },
        "request_uri": {
          "type": "text",
          "index": "not_analyzed"
        }
      }
    }
  }
}
EOF
}`
