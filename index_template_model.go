package elasticsearch

type IndexTemplateMapping struct {
	Properties map[string]interface{} `json:"properties"`
}

type IndexTemplateDoc struct {
	Mappings map[string]IndexTemplateMapping `json:"mappings"`
	Template string                          `json:"template"`
}
