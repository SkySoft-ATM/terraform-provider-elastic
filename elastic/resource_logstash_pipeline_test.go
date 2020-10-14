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
	id := "toto2"
	pipeline := "test pipeline content"
	description := "example description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElasticLogstashDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElasticLogstashPipelineConfigBasic(id, pipeline, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticLogstashPipelineExists("elastic_logstash_pipeline.new"),
					resource.TestCheckResourceAttr("elastic_logstash_pipeline.new", "pipeline", pipeline),
				),
			},
		},
	})
}

func TestAccElasticLogstashPipelineDataSource(t *testing.T) {
	id := "filebeat"
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElasticLogstashPipelineConfigDataSource(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elastic_logstash_pipeline.filebeat", "description", "Pipeline used to consume events from filebeat"),
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

func testAccCheckElasticLogstashPipelineConfigBasic(id, pipeline, description string) string {
	return fmt.Sprintf(`
	resource "elastic_logstash_pipeline" "new" {
		pipeline_id = "%s"
		pipeline 	= "%s"
		description = "%s"
		settings{}
	}
	`, id, pipeline, description)
}

func testAccCheckElasticLogstashPipelineConfigDataSource(id string) string {
	return fmt.Sprintf(`
	data "elastic_logstash_pipeline" "filebeat" {
		pipeline_id = "%s"
	  }
	`, id)
}

func testAccCheckElasticLogstashPipelineExists(n string) resource.TestCheckFunc {
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
		if pipeline.Configuration.Pipeline != pipelineDef {
			return fmt.Errorf("Description does not match: %s", pipeline.Configuration.Pipeline)
		}
		return nil
	}
}
