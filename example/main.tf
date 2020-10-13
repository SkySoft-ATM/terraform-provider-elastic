terraform {
  required_providers {
    elastic = {
      versions = ["0.0.1"]
      source = "hashicorp.com/skysoft-atm/elastic"
    }
  }
}

provider "elastic" {
  kibana_url = "<YOUR_KIBANA_URL>"
  cloud_auth = "<YOUR_CLOUD_AUTH>"
}

data "elastic_logstash_pipeline" "test" {
  pipeline_id = "test"
}

resource "elastic_logstash_pipeline" "fake" {
  pipeline_id = "test"
  pipeline = "Testons donc tout ça"
  description = "Description"
  settings {
    workers = 1
    batch_size = 125
    batch_delay = 50
    queue_checkpoint_writes = 1024
    queue_max_bytes = "1gb"
    queue_type = "persisted"
  }
}

// Output 
// Root Level
output "pipeline_id" {
  value = data.elastic_logstash_pipeline.test.pipeline_id
}

output "description" {
  value = data.elastic_logstash_pipeline.test.description
}

output "pipeline" {
  value = data.elastic_logstash_pipeline.test.pipeline
}

output "username" {
  value = data.elastic_logstash_pipeline.test.username
}

// Settings level
output "settings_workers" {
  value = data.elastic_logstash_pipeline.test.settings
}

output "settings_batch_size" {
  value = data.elastic_logstash_pipeline.test.settings[0].batch_size
}

