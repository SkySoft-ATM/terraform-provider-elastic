terraform {
  required_version = ">= 0.13"
  required_providers {
    elastic = {
      source  = "local/skysoft-atm/elastic"
      version = "0.0.1"
    }
  }
}

resource "elastic_logstash_pipeline" "test" {
  pipeline_id = "fake"
}
