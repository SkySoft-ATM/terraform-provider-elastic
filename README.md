# Elastic Cloud Provider for Terraform
![release](https://github.com/SkySoft-ATM/terraform-provider-elastic/workflows/release/badge.svg)

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
	queue_checkpoint_writes 		= 1024
	queue_max_bytes 			= "1gb"
	queue_type 				= "memory"
  } 
}
```
`pipeline` defintion can be a little be tedious to define inside a JSON, so the `templatefile` [terraform native function](https://www.terraform.io/docs/configuration/functions/templatefile.html) can be used.
Example below illustrates the usage:
```hcl
resource "elastic_logstash_pipeline" "test" {
  pipeline_id = "test"
  pipeline = templatefile("${path.module}/pipeline.conf", {
    CLOUD_ID   = var.cloud_id
    CLOUD_AUTH = var.cloud_auth
  })
  description = "My so great pipeline"
  settings { // Required even if empty (default values will be used)
    	batch_delay				= 50
    	batch_size 				= 125
	workers 				= 1
	queue_checkpoint_writes 		= 1024
	queue_max_bytes 			= "1gb"
	queue_type 				= "memory"
  } 
}
```
An example of `pipeline.conf` is available [here](./example/pipeline.conf)

Using data sources
----------------------
```hcl
data "elastic_logstash_pipeline" "filebeat" {
  pipeline_id = "filebeat"
}
```
