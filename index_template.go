package elasticsearch

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/errwrap"

	v5 "gopkg.in/olivere/elastic.v5"
	v6 "gopkg.in/olivere/elastic.v6"
)

func indexTemplateRead(clients *Clients, id string) ([]byte, error) {
	var template []byte
	switch clients.Version {
	case 6:
		client := clients.V6Client
		resp, err := client.IndexGetTemplate(id).Do(context.Background())
		if err != nil {
			if v6.IsNotFound(err) {
				return nil, nil
			}
			return nil, errwrap.Wrapf("tem[es v6] read index error: {{err}}", err)
		}
		template, err = json.Marshal(resp)
		if err != nil {
			return nil, errwrap.Wrapf("template contains an invalid JSON: {{err}}", err)
		}
	default:
		client := clients.V5Client
		resp, err := client.IndexGetTemplate(id).Do(context.Background())
		if err != nil {
			if v5.IsNotFound(err) {
				return nil, nil
			}
			return nil, errwrap.Wrapf("tem[es v6] read index error: {{err}}", err)
		}
		template, err = json.Marshal(resp)
		if err != nil {
			return nil, errwrap.Wrapf("template contains an invalid JSON: {{err}}", err)
		}
	}
	return template, nil
}

func indexTemplateCreate(clients *Clients, id string, template string, create bool) error {
	switch clients.Version {
	case 6:
		client := clients.V6Client
		_, err := client.IndexPutTemplate(id).BodyString(template).Create(create).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v6] create index error: {{err}}", err)
		}
	default:
		client := clients.V5Client
		_, err := client.IndexPutTemplate(id).BodyString(template).Create(create).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v5] create index error: {{err}}", err)
		}
	}
	return nil
}

func indexTemplateDelete(clients *Clients, id string) error {
	switch clients.Version {
	case 6:
		client := clients.V6Client
		_, err := client.IndexDeleteTemplate(id).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v6] failed to delete index: {{err}}", err)
		}
	default:
		client := clients.V5Client
		_, err := client.IndexDeleteTemplate(id).Do(context.Background())
		if err != nil {
			return errwrap.Wrapf("[es v5] failed to delete index: {{err}}", err)
		}
	}
	return nil
}
