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
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Pipeline name, must be unique.`,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Pipeline definition which will be used by logstash instances.
				Should be composed by 3 sections (input, filter and output).`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Pipeline description.`,
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Token owner used for the pipeline creation.`,
			},
			"settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true, // Workaround, it sounds like the current Kibana API does not behave has expected:
				// If settings is empty, then the API returns null value even if default values are
				// applied...
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      50,
							ValidateFunc: utils.IntAtLeast(1),
							Description: `This setting adjusts the latency of the Logstash pipeline. 
							Pipeline batch delay is the maximum amount of time in milliseconds that 
							Logstash waits for new messages after receiving an event in the current 
							pipeline worker thread.`,
						},
						"batch_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: utils.IntAtLeast(1),
							Description: `This setting defines the maximum number of events an 
							individual worker thread collects before attempting to execute filters 
							and outputs. Larger batch sizes are generally more efficient, but 
							increase memory overhead.`,
						},
						"workers": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: utils.IntAtLeast(1),
							Description: `This setting determines how many threads to run for filter
							 and output processing.`,
						},
						"queue_checkpoint_writes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1024,
							ValidateFunc: utils.IntAtLeast(1),
							Description: `This setting specifies the maximum number of events that
							 may be written to disk before forcing a checkpoint. `,
						},
						"queue_max_bytes": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "1gb",
							Description: `The total capacity of the queue in number of bytes.`,
						},
						"queue_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "memory",
							Description:  `Specify persisted to enable persistent queues.`,
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
	err = c.CreateOrUpdateLogstashPipeline(ctx, &data)
	if err != nil {
		log.Printf("Error : %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(data.ID)

	resourceLogstashPipelineRead(ctx, d, m)

	return diags
}

func resourceLogstashPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	// Warning on errors can be collected in a slice type
	var diags diag.Diagnostics

	pipelineID := d.Id()
	// API crashes if the pipeline_id is not known
	// Let's first look if we can find it in a list
	pipes, err := c.GetLogstashPipelines(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for _, p := range pipes.Pipelines {
		if pipelineID == p.ID {
			found = true
			break
		}
	}

	if found {
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
	}

	return diags
}

func resourceLogstashPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)

	if d.HasChange("description") || d.HasChange("pipeline") || d.HasChange("settings") || d.HasChange("username") {
		data, err := pipelineLogstashData(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = c.CreateOrUpdateLogstashPipeline(ctx, &data)
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

	// Check Prerequisites
	pipelineID := d.Get("pipeline_id").(string)
	if len(pipelineID) == 0 {
		return data, fmt.Errorf("pipeline_id must be defined")
	}

	pipeline := d.Get("pipeline").(string)
	if len(pipeline) == 0 {
		return data, fmt.Errorf("pipeline must be defined")
	}
	// End Check Prerequisites

	data.ID = pipelineID
	config := api.LogstashConfiguration{}
	config.Pipeline = pipeline
	if v, ok := d.GetOk("description"); ok {
		config.Description = v.(string)
	}
	var settings api.Settings
	v := d.Get("settings").([]interface{})
	for _, item := range v {
		i := item.(map[string]interface{})
		settings.PipelineBatchDelay = i["batch_delay"].(int)
		settings.PipelineWorkers = i["workers"].(int)
		settings.PipelineBatchSize = i["batch_size"].(int)
		settings.QueueCheckpointWrites = i["queue_checkpoint_writes"].(int)
		settings.QueueMaxBytes = i["queue_max_bytes"].(string)
		settings.QueueType = i["queue_type"].(string)
	}
	config.Settings = &settings

	data.Configuration = &config
	return data, nil
}
