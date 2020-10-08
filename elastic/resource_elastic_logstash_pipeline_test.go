package elastic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/skysoft-atm/terraform-provider-elastic/api"
)

func TestAccLogstashPipeline_basic(t *testing.T) {
	var pipelineDef api.LogstashPipeline
	id := "fake"
	pipeline := "test pipeline content"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogstashDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogstashPipelineConfigBasic(id, pipeline),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogstashPipelineExists("logstash_pipeline.test", &pipelineDef),
					testAccCheckLogstashPipelineAttributes(&pipelineDef, pipeline),
					resource.TestCheckResourceAttr("logstash_pipeline.test", "pipeline", pipeline),
				),
			},
		},
	})
}

func testAccCheckLogstashDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "logstash_pipeline" {
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

func testAccCheckLogstashPipelineConfigBasic(id, description string) string {
	return fmt.Sprintf(`
	resource "logstash_pipeline" "test" {
		pipeline_id = "%s"
		description = "%s"
	}
	`, id, description)
}

func testAccCheckLogstashPipelineExists(n string, pipeline *api.LogstashPipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*api.Client)

		foundTask, err := client.GetLogstashPipeline(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if foundTask.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*pipeline = *foundTask
		return nil
	}
}

func testAccCheckLogstashPipelineAttributes(pipeline *api.LogstashPipeline, pipelineDef string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pipeline.Pipeline != pipelineDef {
			return fmt.Errorf("Description does not match: %s", pipeline.Pipeline)
		}
		return nil
	}
}
