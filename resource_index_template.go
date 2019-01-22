package elasticsearch

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceElasticsearchIndexTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticsearchIndexTemplateCreate,
		Read:   resourceElasticsearchIndexTemplateRead,
		Update: resourceElasticsearchIndexTemplateUpdate,
		Delete: resourceElasticsearchIndexTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) (json string) {
					json, _ = normalizeJsonString(v)
					return
				},
			},
		},
	}
}

func resourceElasticsearchIndexTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	if err := indexTemplateCreate(meta.(*Clients), name, d.Get("template").(string), true); err != nil {
		return err
	}

	d.SetId(name)
	return resourceElasticsearchIndexTemplateRead(d, meta)
}

func resourceElasticsearchIndexTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	if err := indexTemplateCreate(meta.(*Clients), name, d.Get("template").(string), false); err != nil {
		return err
	}

	d.SetId(name)
	return resourceElasticsearchIndexTemplateRead(d, meta)
}

func resourceElasticsearchIndexTemplateRead(d *schema.ResourceData, meta interface{}) error {
	template, err := indexTemplateRead(meta.(*Clients), d.Id())
	if err != nil {
		return err
	}

	d.Set("template", template)
	d.Set("name", d.Id())

	return nil
}

func resourceElasticsearchIndexTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	if err := indexTemplateDelete(meta.(*Clients), d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
