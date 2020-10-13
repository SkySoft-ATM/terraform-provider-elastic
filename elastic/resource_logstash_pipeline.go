package elastic

import (
	"context"
	"fmt"
	"log"

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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
						},
						"batch_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
						},
						"workers": {
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

	data, err := pipelineLogstashData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("YOUHOU : %v", data)
	err = c.CreateOrUpdateLogstashPipeline(ctx, &data, data.ID)
	if err != nil {
		log.Printf("Error : %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId(data.ID)
	log.Printf("Everything is good")
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
		if err := d.Set(key, value); err != nil {
			diag.FromErr(err)
		}
	}

	return diags
}

func resourceLogstashPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	id := d.Id()

	if d.HasChange("description") || d.HasChange("pipeline") || d.HasChange("username") || d.HasChange("settings") {
		data, err := pipelineLogstashData(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = c.CreateOrUpdateLogstashPipeline(ctx, &data, id)
		if err != nil {
			return diag.FromErr(err)
		}
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

func pipelineLogstashData(d *schema.ResourceData) (api.LogstashPipeline, error) {
	data := api.LogstashPipeline{}
	pipelineID := d.Get("pipeline_id").(string)
	if len(pipelineID) == 0 {
		return data, fmt.Errorf("pipeline_id must be defined")
	}
	data.ID = pipelineID

	if v, ok := d.GetOk("description"); ok {
		data.Description = v.(string)
	}

	if v, ok := d.GetOk("pipeline"); ok {
		data.Pipeline = v.(string)
	}

	settings := api.Settings{}
	vSettings := d.Get("settings").([]interface{})
	for _, item := range vSettings {
		i := item.(map[string]interface{})
		settings.PipelineBatchDelay = i["batch_delay"].(int)
		settings.PipelineWorkers = i["workers"].(int)
		settings.PipelineBatchSize = i["batch_size"].(int)
		settings.QueueCheckpointWrites = i["queue_checkpoint_writes"].(int)
		settings.QueueMaxBytes = i["queue_max_bytes"].(string)
		settings.QueueType = i["queue_type"].(string)
	}
	data.Settings = &settings
	return data, nil
}
