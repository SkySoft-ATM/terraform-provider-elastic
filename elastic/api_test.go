package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	refPipeline = []byte(`{
		"id": "test",
		"description": "Used to test terraform provider",
		"username": "chambodn@skysoft-atm.com",
		"pipeline": "input {\n    beats {\n        port => 5044\n    }\n}\nfilter {\n    if \"eu.gcr.io/sk-private-registry/skysoft-atm/\" not in [kubernetes][container][image] and \"elasticsearch\" not in [kubernetes][labels][name] {\n        drop {}\n    }\n    if [kubernetes][labels][name] == \"elasticsearch\" {\n        grok {\n            match => [ \"message\", \"[%{TIMESTAMP_ISO8601:timestamp}][%{DATA:level}%{SPACE}][%{DATA:source}%{SPACE}]%{SPACE}%{GREEDYDATA:message}\" ]\n            overwrite => [ \"message\" ]\n        }\n    }\n    if [kubernetes][labels][name] == \"ems\" {\n        grok {\n            match => { \"message\" => \"%{TIMESTAMP_ISO8601:logdate} %{LOGLEVEL:level}: %{GREEDYDATA:message}\" }\n            overwrite => [\"message\"]\n        \n        date {\n            match => [ \"logdate\", \"ISO8601\", \"yyyy-MM-dd HH:mm:ss,SSS\", \"yyyy-MM-dd HH:mm:ss.SSS\" ]\n            remove_field => [ \"logdate\" ]\n        }\n    }\n    date {\n        match => [ \"timestamp\", \"ISO8601\", \"yyyy-MM-dd HH:mm:ss,SSS\", \"yyyy-MM-dd HH:mm:ss.SSS\" ]\n    }\n}\noutput {\n    elasticsearch {\n        index => \"filebeat-%{+yyyy.MM.dd}\"\n        cloud_id => \"Test\"\n        cloud_auth => \"Test\"\n    }\n}",
		"settings": {
			"pipeline.batch.delay": 50,
			"pipeline.batch.size": 125,
			"pipeline.workers": 1,
			"queue.checkpoint.writes": 1024,
			"queue.max_bytes": "1gb",
			"queue.type": "persistent"
		}
	}`)
	client      *Client
	pipelineRef LogstashPipeline
)

func init() {
	client = NewClient(os.Getenv("CLOUD_AUTH"), "https://80128f85e27c4ed8bab925c5cc6b811c.europe-west1.gcp.cloud.es.io:9243")
	err := json.Unmarshal(refPipeline, &pipelineRef)
	if err != nil {
		fmt.Errorf("error trying to load pipeline definition %s", err.Error())
	}
}

func TestGetPipeline(t *testing.T) {
	res, err := client.GetLogstashPipeline(context.Background(), "test")

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")

	assert.Equal(t, pipelineRef.ID, res.ID, "expecting same IDs")
	assert.Equal(t, pipelineRef.Description, res.Description, "expecting same description")
	assert.Equal(t, pipelineRef.Username, res.Username, "expecting same username")
	assert.Equal(t, pipelineRef.Settings, res.Settings, "expecting same settings")
	assert.Equal(t, pipelineRef.Pipeline, res.Pipeline, "expecting same pipeline definition")
}

func TestCreatePipeline(t *testing.T) {
	settings := NewLogstashPipelineSettings(50, 125, 1, 1024, "1gb", "memory")
	logstashPipeline := NewLogstashPipeline("", "Used to test terraform provider", "Test pipeline", settings)
	err := client.CreateLogstashPipeline(context.Background(), logstashPipeline, "empty")

	assert.Nil(t, err, "expecting nil error")
}

func TestUpdatePipeline(t *testing.T) {
	settings := NewLogstashPipelineSettings(50, 125, 1, 1024, "1gb", "memory")
	logstashPipeline := NewLogstashPipeline("", "Used to test terraform provider", "Test pipeline", settings)
	err := client.UpdateLogstashPipeline(context.Background(), logstashPipeline, "empty")

	assert.Nil(t, err, "expecting nil error")
}

func TestDeletePipeline(t *testing.T) {
	err := client.DeleteLogstashPipeline(context.Background(), "empty")
	assert.Nil(t, err, "expecting nil error")
}
