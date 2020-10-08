package elastic

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
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
			"pipeline_batch_delay": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pipeline_batch_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pipeline_workers": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"queue_checkpoint_writes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"queue_max_bytes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateValueFunc([]string{"memory", "persisted"}),
			},
		},
		Create: resourceLogstashPipelineCreate,
		Read:   resourceLogstashPipelineRead,
		Update: resourceLogstashPipelineUpdate,
		Delete: resourceLogstashPipelineDelete,
	}
}

func resourceLogstashPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("pipeline_id").(string)
	data := pipelineLogstashData(d, meta)

	client := meta.(*api.Client)
	err := client.CreateOrUpdateLogstashPipeline(context.Background(), &data, id)
	if err != nil {
		return err
	}
	d.SetId(id)
	return resourceLogstashPipelineRead(d, meta)
}

func resourceLogstashPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	pipelineID := d.Id()

	var pipeline *api.LogstashPipeline
	pipeline, err := client.GetLogstashPipeline(context.Background(), pipelineID)
	if err != nil {
		return err
	}
	d.Set("description", pipeline.Description)
	d.Set("pipeline", pipeline.Pipeline)
	d.Set("username", pipeline.Username)
	d.Set("pipeline_batch_delay", pipeline.Settings.PipelineBatchDelay)
	d.Set("pipeline_batch_size", pipeline.Settings.PipelineBatchSize)
	d.Set("pipeline_workers", pipeline.Settings.PipelineWorkers)
	d.Set("queue_checkpoint_writes", pipeline.Settings.QueueCheckpointWrites)
	d.Set("queue_max_bytes", pipeline.Settings.QueueMaxBytes)
	d.Set("queue_type", pipeline.Settings.QueueType)
	return nil
}

func resourceLogstashPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	pipelineID := d.Id()
	err := client.DeleteLogstashPipeline(context.Background(), pipelineID)
	if err != nil {
		return err
	}
	return nil
}

func resourceLogstashPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	id := d.Id()

	if d.HasChange("description") || d.HasChange("pipeline") || d.HasChange("username") || d.HasChange("pipeline_batch_delay") ||
		d.HasChange("pipeline_batch_size") || d.HasChange("pipeline_workers") || d.HasChange("queue_checkpoint_writes") ||
		d.HasChange("queue_max_bytes") || d.HasChange("queue_type") {
		data := pipelineLogstashData(d, meta)
		err := client.CreateOrUpdateLogstashPipeline(context.Background(), &data, id)
		if err != nil {
			return nil
		}

	}
	return nil
}

func pipelineLogstashData(d *schema.ResourceData, meta interface{}) api.LogstashPipeline {
	data := api.LogstashPipeline{}
	if v, ok := d.GetOk("description"); ok {
		data.Description = v.(string)
	}

	if v, ok := d.GetOk("pipeline"); ok {
		data.Pipeline = v.(string)
	}

	settings := api.Settings{}
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
