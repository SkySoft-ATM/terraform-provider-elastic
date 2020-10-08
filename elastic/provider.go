package elastic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
)

// Provider is used by terraform to instantiate Provider object
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cloud_auth": {
				Type:        schema.TypeString,
				Description: "Your CLOUD_AUTH credentials",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_AUTH", ""),
			},
			"kibana_url": {
				Type:        schema.TypeString,
				Description: "Kibana URL",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_URL", ""),
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
		client := api.NewClient(d.Get("cloud_auth").(string), d.Get("kibana_url").(string))
		return client, nil
	}
}
