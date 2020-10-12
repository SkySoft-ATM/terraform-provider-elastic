package elastic

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
	"github.com/skysoft-atm/terraform-provider-elastic/utils"
)

func resourceLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_batch_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
						},
						"pipeline_batch_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
						},
						"pipeline_workers": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
						},
						"queue_checkpoint_writes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
						},
						"queue_max_bytes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: utils.StringInSlice([]string{"memory", "persisted"}, false),
						},
					},
				},
			},
		},
		CreateContext: resourceLogstashPipelineCreate,
		ReadContext:   resourceLogstashPipelineRead,
		UpdateContext: resourceLogstashPipelineUpdate,
		DeleteContext: resourceLogstashPipelineDelete,
	}
}

func resourceLogstashPipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	// Warning on errors can be collected in a slice type
	var diags diag.Diagnostics

	data := pipelineLogstashData(d)

	err := c.CreateOrUpdateLogstashPipeline(ctx, &data, data.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(data.ID)

	return diags
}

func resourceLogstashPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	// Warning on errors can be collected in a slice type
	var diags diag.Diagnostics

	pipelineID := d.Id()

	pipeline, err := c.GetLogstashPipeline(ctx, pipelineID)
	if err != nil {
		return diag.FromErr(err)
	}

	pl := flattenLogstashPipelineData(pipeline)
	for key, value := range pl {
		d.Set(key, value)
	}

	return diags
}

func resourceLogstashPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	id := d.Id()

	if d.HasChange("description") || d.HasChange("pipeline") || d.HasChange("username") || d.HasChange("pipeline_batch_delay") ||
		d.HasChange("pipeline_batch_size") || d.HasChange("pipeline_workers") || d.HasChange("queue_checkpoint_writes") ||
		d.HasChange("queue_max_bytes") || d.HasChange("queue_type") {
		data := pipelineLogstashData(d)
		err := c.CreateOrUpdateLogstashPipeline(ctx, &data, id)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourceLogstashPipelineRead(ctx, d, m)
}

func resourceLogstashPipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pipelineID := d.Id()

	err := c.DeleteLogstashPipeline(ctx, pipelineID)
	if err != nil {
		return diag.FromErr(err)
	}
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func pipelineLogstashData(d *schema.ResourceData) api.LogstashPipeline {
	data := api.LogstashPipeline{}
	if v, ok := d.GetOk("description"); ok {
		data.Description = v.(string)
	}

	if v, ok := d.GetOk("pipeline"); ok {
		data.Pipeline = v.(string)
	}

	settings := api.Settings{}
	// if v, ok := d.GetOk("settings"); ok {

	// }
	if v, ok := d.GetOk("pipeline_batch_delay"); ok {
		settings.PipelineBatchDelay = v.(int)
	}

	if v, ok := d.GetOk("pipeline_batch_size"); ok {
		settings.PipelineBatchSize = v.(int)
	}

	if v, ok := d.GetOk("pipeline_workers"); ok {
		settings.PipelineWorkers = v.(int)
	}

	if v, ok := d.GetOk("queue_checkpoint_writes"); ok {
		settings.QueueCheckpointWrites = v.(int)
	}

	if v, ok := d.GetOk("queue_max_bytes"); ok {
		settings.QueueMaxBytes = v.(string)
	}

	if v, ok := d.GetOk("queue_type"); ok {
		settings.QueueType = v.(string)
	}
	data.Settings = &settings
	return data
}
