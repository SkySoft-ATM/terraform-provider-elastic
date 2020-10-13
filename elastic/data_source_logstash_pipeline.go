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
			"pipeline_id": {
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"batch_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"workers": {
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
	if err != nil {
		return diag.FromErr(err)
	}

	pl := flattenLogstashPipelineData(pipeline)
	for key, value := range pl {
		if err := d.Set(key, value); err != nil {
			diag.FromErr(err)
		}
	}

	d.SetId(pipeline.ID)

	return diags
}

func flattenLogstashPipelineData(pipeline *api.LogstashPipeline) map[string]interface{} {
	lp := make(map[string]interface{})
	if pipeline != nil {
		lp["pipeline_id"] = pipeline.ID
		lp["description"] = pipeline.Description
		lp["username"] = pipeline.Username
		lp["pipeline"] = pipeline.Pipeline
		lp["settings"] = flattenSettings(pipeline.Settings)
	}
	return lp
}

func flattenSettings(settings *api.Settings) []interface{} {
	s := make(map[string]interface{})
	s["workers"] = settings.PipelineWorkers
	s["batch_size"] = settings.PipelineBatchSize
	s["batch_delay"] = settings.PipelineBatchDelay
	s["queue_checkpoint_writes"] = settings.QueueCheckpointWrites
	s["queue_max_bytes"] = settings.QueueMaxBytes
	s["queue_type"] = settings.QueueType
	return []interface{}{s}
}
