package elastic

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"settings": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
			},
		},
		Create: resourceLogstashPipelineCreate,
		Read:   resourceLogstashPipelineRead,
		Update: resourceLogstashPipelineUpdate,
		Delete: resourceLogstashPipelineDelete,
	}
}

func resourceLogstashPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLogstashPipelineRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLogstashPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLogstashPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
