package elastic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider is used by terraform to instantiate Provider object
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Your elastic API key",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELASTIC_CLOUD", nil),
			},
			"kibana_host": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Kibana host",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_HOST", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"logstash_pipeline": resourceLogstashPipeline(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"logstash_pipeline": dataSourceLogstashPipeline(),
		},
		ConfigureFunc: configureFunc(),
	}
}

func configureFunc() func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		return nil, nil
	}
}
