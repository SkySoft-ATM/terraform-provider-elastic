package elastic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
)

func TestAccElasticLogstashPipeline_basic(t *testing.T) {
	var pipelineDef api.LogstashPipeline
	id := "fake"
	pipeline := "test pipeline content"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElasticLogstashDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElasticLogstashPipelineConfigBasic(id, pipeline),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticLogstashPipelineExists("elastic_logstash_pipeline.new", &pipelineDef),
					//testAccCheckElasticLogstashPipelineAttributes(&pipelineDef, pipeline),
					//resource.TestCheckResourceAttr("logstash_pipeline.test", "pipeline", pipeline),
				),
			},
		},
	})
}

func testAccCheckElasticLogstashDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elastic_logstash_pipeline" {
			continue
		}

		// Try to find the task
		_, err := client.GetLogstashPipeline(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Task still exists")
		}

	}
	return nil
}

func testAccCheckElasticLogstashPipelineConfigBasic(id, description string) string {
	return fmt.Sprintf(`
	resource "elastic_logstash_pipeline" "new" {
		pipeline_id = "%s"
		description = "%s"
	}
	`, id, description)
}

func testAccCheckElasticLogstashPipelineExists(n string, pipeline *api.LogstashPipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		return nil
	}
}

func testAccCheckElasticLogstashPipelineAttributes(pipeline *api.LogstashPipeline, pipelineDef string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pipeline.Pipeline != pipelineDef {
			return fmt.Errorf("Description does not match: %s", pipeline.Pipeline)
		}
		return nil
	}
}
