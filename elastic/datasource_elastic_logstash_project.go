package elastic

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
)

func dataSourceLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Read: dataSourceLogstashPipelineRead,
	}
}

func dataSourceLogstashPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	id := d.Get("pipeline_id").(string)
	pipeline, err := client.GetLogstashPipeline(context.Background(), id)
	if err != nil {
		return err
	}
	d.SetId(pipeline.ID)
	return nil
}
