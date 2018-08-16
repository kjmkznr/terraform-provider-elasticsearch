Terraform provider for elasticsearch
====================================

A [Terraform](https://www.terraform.io/) plugin that provides resources for elasticsearch

[![Build Status](https://travis-ci.org/kjmkznr/terraform-provider-elasticsearch.svg?branch=master)](https://travis-ci.org/kjmkznr/terraform-provider-elasticsearch)

Install
-------

* Download the latest release for your platform.
* Rename the executable to `terraform-provider-elasticsearch`

Provider Configuration
----------------------

### Example

```
provider "elasticsearch" {
  url        = "http://localhost:9200"
  es_version = "v5"
}
```

or

```
provider "elasticsearch" {}
```

### Reference

* `url` - (Required) Endpoint of elasticsearch
* `es_version` - (Default v5) Version of elasticsearch

Resources
---------

### `elasticsearch_index_template`

Configure a Elasticsearch index template.

#### Example
```
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
}
```

Build
-----

```
$ make build
```

Testing
-------

```
$ ELASTICSEARCH_URL="http://127.0.0.1:9200" ELASTICSEARCH_VERSION=v5 make testacc
```

Licence
-------

Mozilla Public License, version 2.0

Author
------

[KOJIMA Kazunori](https://github.com/kjmkznr)

