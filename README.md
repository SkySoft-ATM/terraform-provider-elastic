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

