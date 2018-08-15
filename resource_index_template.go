package elasticsearch

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"

	v5 "gopkg.in/olivere/elastic.v5"
	v6 "gopkg.in/olivere/elastic.v6"
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
	clients := meta.(*Clients)
	name := d.Get("name").(string)
	template := d.Get("template").(string)

	if clients.Version == 6 {
		_, err := clients.V6Client.IndexPutTemplate(name).BodyString(template).Create(true).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v6] create index error: {{err}}", err)
		}
	} else {
		_, err := clients.V5Client.IndexPutTemplate(name).BodyString(template).Create(true).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v5] create index error: {{err}}", err)
		}
	}

	d.SetId(d.Get("name").(string))
	return resourceElasticsearchIndexTemplateRead(d, meta)
}

func resourceElasticsearchIndexTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*Clients)
	name := d.Get("name").(string)
	template := d.Get("template").(string)

	if clients.Version == 6 {
		_, err := clients.V6Client.IndexPutTemplate(name).BodyString(template).Create(false).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v6] update index error: {{err}}", err)
		}
	} else {
		//client := clients.V5Client
		_, err := clients.V5Client.IndexPutTemplate(name).BodyString(template).Create(false).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v5] update index error: {{err}}", err)
		}
	}

	d.SetId(d.Get("name").(string))
	return resourceElasticsearchIndexTemplateRead(d, meta)
}

func resourceElasticsearchIndexTemplateRead(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*Clients)
	if clients.Version == 6 {
		resp, err := clients.V6Client.IndexGetTemplate(d.Id()).Do(context.Background())
		if err != nil {
			if v6.IsNotFound(err) {
				d.Set("name", "")
				return nil
			}
			return errwrap.Wrapf("[es v6] read index error: {{err}}", err)
		}

		template, err := json.Marshal(resp)
		if err != nil {
			return errwrap.Wrapf("template contains an invalid JSON: {{err}}", err)
		}
		d.Set("template", template)

	} else {
		resp, err := clients.V5Client.IndexGetTemplate(d.Id()).Do(context.Background())
		if err != nil {
			if v5.IsNotFound(err) {
				d.Set("name", "")
				return nil
			}
			return errwrap.Wrapf("[es v5] read index error: {{err}}", err)
		}

		template, err := json.Marshal(resp)
		if err != nil {
			return errwrap.Wrapf("template contains an invalid JSON: {{err}}", err)
		}
		d.Set("template", template)

	}

	d.Set("name", d.Id())

	return nil
}

func resourceElasticsearchIndexTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*Clients)
	//var err error
	if clients.Version == 6 {
		_, err := clients.V6Client.IndexDeleteTemplate(d.Id()).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v6] failed to delete index: {{err}}", err)
		}

	} else {
		_, err := clients.V5Client.IndexDeleteTemplate(d.Id()).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v5] failed to delete index: {{err}}", err)
		}
	}

	d.SetId("")

	return nil
}
