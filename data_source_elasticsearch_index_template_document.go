package elasticsearch

import (
	"encoding/json"

	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceElasticsearchIndexTemplateDocument() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElasticsearchIndexTemplateDocumentRead,

		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"property": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"index": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceElasticsearchIndexTemplateDocumentRead(d *schema.ResourceData, meta interface{}) error {
	doc := &IndexTemplateDoc{}

	if template, ok := d.GetOk("template"); ok {
		doc.Template = template.(string)
	}

	var cfgMappings = d.Get("mapping").([]interface{})
	mappings := make(map[string]IndexTemplateMapping)
	for _, mappingI := range cfgMappings {
		cfgMap := mappingI.(map[string]interface{})
		mapping := IndexTemplateMapping{}

		if cfgProperties := cfgMap["property"].(*schema.Set).List(); len(cfgProperties) > 0 {
			properties := make(map[string]interface{})
			for _, propI := range cfgProperties {
				cfgProp := propI.(map[string]interface{})
				prop := map[string]string{
					"type": cfgProp["type"].(string),
				}
				if f := cfgProp["format"].(string); f != "" {
					prop["format"] = cfgProp["format"].(string)
				}
				if f := cfgProp["index"].(string); f != "" {
					prop["index"] = cfgProp["index"].(string)
				}

				key := cfgProp["name"].(string)
				properties[key] = prop
			}
			mapping.Properties = properties
		}

		key := cfgMap["type"].(string)
		mappings[key] = mapping
	}
	doc.Mappings = mappings

	jsonDoc, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)

	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}
