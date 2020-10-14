terraform {
  required_providers {
    elastic = {
      versions = ["0.0.1"]
      source   = "hashicorp.com/skysoft-atm/elastic"
    }
  }
}
/*
  Variables
*/
variable "kibana_url" {
  type = string
}

variable "cloud_auth" {
  type = string
}


provider "elastic" {
  kibana_url = var.kibana_url
  cloud_auth = var.cloud_auth
}

data "elastic_logstash_pipeline" "filebeat" {
  pipeline_id = "filebeat"
}

resource "elastic_logstash_pipeline" "test" {
  pipeline_id = "test"
  pipeline    = "Testons donc tout ça"
  description = "Description"
  settings {}
}

// Output Filebeat
// Root Level
output "filebeat_pipeline_id" {
  value = data.elastic_logstash_pipeline.filebeat.pipeline_id
}

output "filebeat_description" {
  value = data.elastic_logstash_pipeline.filebeat.description
}

output "filebeat_pipeline" {
  value = data.elastic_logstash_pipeline.filebeat.pipeline
}

output "filebeat_username" {
  value = data.elastic_logstash_pipeline.filebeat.username
}

// Settings level
output "filebeat_settings_workers" {
  value = data.elastic_logstash_pipeline.filebeat.settings[0].workers
}

output "filebeat_settings_batch_size" {
  value = data.elastic_logstash_pipeline.filebeat.settings[0].batch_size
}

// Output test
// Root Level
output "test_pipeline_id" {
  value = elastic_logstash_pipeline.test.pipeline_id
}

output "test_description" {
  value = elastic_logstash_pipeline.test.description
}

output "test_pipeline" {
  value = elastic_logstash_pipeline.test.pipeline
}

output "test_username" {
  value = elastic_logstash_pipeline.test.username
}

// Settings level
output "test_settings_workers" {
  value = elastic_logstash_pipeline.test.settings[0].workers
}

output "test_settings_batch_size" {
  value = elastic_logstash_pipeline.test.settings[0].batch_size
}
