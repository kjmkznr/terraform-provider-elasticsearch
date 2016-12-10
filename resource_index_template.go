package elasticsearch

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	elastic "gopkg.in/olivere/elastic.v5"
)

func resourceElasticsearchIndexTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticsearchIndexTemplateCreate,
		Read:   resourceElasticsearchIndexTemplateRead,
		Update: resourceElasticsearchIndexTemplateCreate,
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
	svc := elastic.NewIndicesPutTemplateService(meta.(*elastic.Client))
	svc.Name(d.Get("name").(string))
	svc.BodyString(d.Get("template").(string))

	_, err := svc.Do(context.Background())
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceElasticsearchIndexTemplateRead(d, meta)
}

func resourceElasticsearchIndexTemplateRead(d *schema.ResourceData, meta interface{}) error {
	svc := elastic.NewIndicesGetTemplateService(meta.(*elastic.Client))
	svc.Name(d.Id())

	resp, err := svc.Do(context.Background())
	if err != nil {
		return err
	}

	d.Set("name", d.Id())

	template, err := json.Marshal(resp)
	if err != nil {
		return errwrap.Wrapf("template contains an invalid JSON: {{err}}", err)
	}
	d.Set("template", template)

	return nil
}

func resourceElasticsearchIndexTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	svc := elastic.NewIndicesDeleteTemplateService(meta.(*elastic.Client))
	svc.Name(d.Id())

	_, err := svc.Do(context.Background())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
