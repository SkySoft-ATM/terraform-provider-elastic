package elastic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
)

// Provider is used by terraform to instantiate Provider object
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cloud_auth": {
				Type:        schema.TypeString,
				Description: "Your CLOUD_AUTH credentials",
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD_AUTH", nil),
			},
			"kibana_url": {
				Type:        schema.TypeString,
				Description: "Kibana URL",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"elastic_logstash_pipeline": resourceLogstashPipeline(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"elastic_logstash_pipeline": dataSourceLogstashPipeline(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cloudAuth := d.Get("cloud_auth").(string)
	kibanaURL := d.Get("kibana_url").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (cloudAuth != "") && (kibanaURL != "") {
		c := api.NewClient(cloudAuth, kibanaURL)
		return c, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to create Logstash client",
		Detail:   "KIBANA_URL or CLOUD_AUTH are not specified",
	})

	return nil, diags
}
