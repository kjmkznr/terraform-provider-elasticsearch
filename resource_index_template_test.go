package elasticsearch

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIndexTemplate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndexTemplateDestroy,
		Steps: []resource.TestStep{
			{
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
			{
				Config: testAccCheckIndexDynamicTemplateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foo", "name", "dynamic1"),
					resource.TestMatchResourceAttr(
						"elasticsearch_index_template.foo", "template",
						regexp.MustCompile(`"integers":{`)),
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
			{
				Config: testAccCheckIndexTemplateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foobar", "name", "template1"),
				),
			},
			{
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
			{
				Config: testAccCheckIndexDynamicTemplateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foo", "name", "dynamic1"),
				),
			},
			{
				Config: testAccCheckIndexDynamicTemplateConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"elasticsearch_index_template.foo", "name", "dynamic1"),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foo", "template",
						regexp.MustCompile(`"index_patterns":\["dynamic-\*`)),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foo", "template",
						regexp.MustCompile(`"fields":{`)),
					resource.TestMatchResourceAttr("elasticsearch_index_template.foo", "template",
						regexp.MustCompile(`"strings":{`)),
				),
			},
		},
	})
}

func testAccCheckIndexTemplateDestroy(s *terraform.State) error {
	clients := testAccProvider.Meta().(*Clients)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticsearch_index_template" {
			continue
		}

		err := indexTemplateDelete(clients, rs.Primary.ID)
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
          "index": "false"
        }
      }
    }
  }
}
EOF
}`

const testAccCheckIndexDynamicTemplateConfig_basic = `
resource "elasticsearch_index_template" "foo" {
    name         = "dynamic1"
	template = <<EOF
{
	"index_patterns": ["dynamic-*"],
	"settings": {
		"number_of_shards": 1
	},
	"mappings": {
		"_doc": {
			"dynamic_templates": [
				{
					"integers": {
						"match_mapping_type": "long",
						"mapping": {
							"type": "integer"
						}
					}
				},
				{
					"hoge": {
						"match_mapping_type": "long",
						"mapping": {
							"type": "integer"
						}
					}
				}
			]
		}
	}
}
EOF
}`

const testAccCheckIndexDynamicTemplateConfig_update = `
resource "elasticsearch_index_template" "foo" {
    name         = "dynamic1"
	template = <<EOF
{
	"index_patterns": ["dynamic-*"],
	"settings": {
		"number_of_shards": 1
	},
	"mappings": {
		"_doc": {
			"dynamic_templates": [
				{
					"integers": {
						"match_mapping_type": "long",
						"mapping": {
							"type": "integer"
						}
					}
				},
				{
					"strings": {
						"match_mapping_type": "string",
						"mapping": {
							"type": "text",
							"fields": {
								"raw": {
									"type":  "keyword",
									"ignore_above": 256
								}
							}
						}
					}
				}
			]
		}
	}
}
EOF
}`
