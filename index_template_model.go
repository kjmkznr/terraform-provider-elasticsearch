package elasticsearch

type IndexTemplateMapping struct {
	Properties map[string]interface{} `json:"properties"`
	Source     map[string]interface{} `json:"_source,omitempty"`
	All        map[string]interface{} `json:"_all,omitempty"`
}

type IndexTemplateSetting struct {
	NumberOfShards int `json:"number_of_shards"`
}

type IndexTemplateDoc struct {
	Settings IndexTemplateSetting            `json:"settings"`
	Mappings map[string]IndexTemplateMapping `json:"mappings"`
	Template string                          `json:"template"`
}
