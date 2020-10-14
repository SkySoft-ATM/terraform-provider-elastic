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
				Type:        schema.TypeString,
				Required:    true,
				Description: `Pipeline name, must be unique.`,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Pipeline definition which will be used by logstash instances.
				Should be composed by 3 sections (input, filter and output).`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Pipeline description.`,
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Token owner used for the pipeline creation.`,
			},
			"settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_delay": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `This setting adjusts the latency of the Logstash pipeline. 
							Pipeline batch delay is the maximum amount of time in milliseconds that 
							Logstash waits for new messages after receiving an event in the current 
							pipeline worker thread.`,
						},
						"batch_size": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `This setting defines the maximum number of events an 
							individual worker thread collects before attempting to execute filters 
							and outputs. Larger batch sizes are generally more efficient, but 
							increase memory overhead.`,
						},
						"workers": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `This setting determines how many threads to run for filter
							and output processing.`,
						},
						"queue_checkpoint_writes": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: `This setting specifies the maximum number of events that
							may be written to disk before forcing a checkpoint. `,
						},
						"queue_max_bytes": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The total capacity of the queue in number of bytes.`,
						},
						"queue_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Persistent mode for queues (persisted or memory).`,
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

	// API crashes if the pipeline_id is not known
	// Let's first look if we can find it in a list
	pipes, err := c.GetLogstashPipelines(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for _, p := range pipes.Pipelines {
		if id == p.ID {
			found = true
			break
		}
	}

	if found {
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
	}
	d.SetId(id)

	return diags
}

func flattenLogstashPipelineData(pipeline *api.LogstashPipeline) map[string]interface{} {
	lp := make(map[string]interface{})
	if pipeline != nil {
		lp["pipeline_id"] = pipeline.ID
		lp["description"] = pipeline.Configuration.Description
		lp["username"] = pipeline.Configuration.Username
		lp["pipeline"] = pipeline.Configuration.Pipeline
		lp["settings"] = flattenSettings(pipeline.Configuration.Settings)
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
