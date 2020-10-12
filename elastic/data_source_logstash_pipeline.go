package elastic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
)

func dataSourceLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogstashPipelineRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"settings": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_batch_delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pipeline_batch_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pipeline_workers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queue_checkpoint_writes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queue_max_bytes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLogstashPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("pipeline_id").(string)
	pipeline, err := c.GetLogstashPipeline(ctx, id)
	pl := flattenLogstashPipelineData(pipeline)
	for key, value := range pl {
		d.Set(key, value)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pipeline.ID)

	return diags
}

func flattenLogstashPipelineData(pipeline *api.LogstashPipeline) map[string]interface{} {
	lp := make(map[string]interface{})
	if pipeline != nil {

		lp["pipeline_id"] = pipeline.ID
		lp["pipeline_description"] = pipeline.Description
		lp["pipeline_username"] = pipeline.Username
		lp["pipeline"] = pipeline.Pipeline
		lp["pipeline_workers"] = pipeline.Settings.PipelineWorkers
		lp["pipeline_batch_size"] = pipeline.Settings.PipelineBatchSize
		lp["pipeline_batch_delay"] = pipeline.Settings.PipelineBatchDelay
		lp["queue_checkpoint_writes"] = pipeline.Settings.QueueCheckpointWrites
		lp["queue_max_bytes"] = pipeline.Settings.QueueMaxBytes
		lp["queue_type"] = pipeline.Settings.QueueType
	}
	return lp
}
