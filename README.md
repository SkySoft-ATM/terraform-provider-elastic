# Elastic Cloud Provider for Terraform

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.13+

Using the provider
----------------------
```hcl
provider "elastic" {
  kibana_url = var.kibana_url
  cloud_auth = var.cloud_auth
}
```
where `kibana_url` is the Kibana URL exposing logstash pipeline API and `cloud_auth` the credential to authenticate on kibana api (please note that at this stage only Basic Authentication is supported and provider should not be configured with identity managed externally)

Upgrading the provider
----------------------

The elastic provider doesn't upgrade automatically once you've started using it. After a new release you can run

```bash
terraform init -upgrade
```
to upgrade to the latest stable version of the elastic provider. 

Creating pipeline resources
----------------------
```hcl
resource "elastic_logstash_pipeline" "test" {
  pipeline_id = "test"
  pipeline = "input { stdin {} } output { stdout {} }"
  description = "My so great pipeline"
  settings { // Required even if empty (default values will be used)
	batch_delay				= 50
    batch_size 				= 125
	workers 				= 1
	queue_checkpoint_writes = 1024
	queue_max_bytes 		= "1gb"
	queue_type 				= "memory"
  } 
}
```

Using data sources
----------------------
```hcl
data "elastic_logstash_pipeline" "filebeat" {
  pipeline_id = "filebeat"
}
```
