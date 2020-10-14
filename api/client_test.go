package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/lithammer/shortuuid/v3"
	"github.com/stretchr/testify/assert"
)

var (
	refPipeline = []byte(`{
		"description": "Used to test terraform provider",
		"pipeline": "input {\n    beats {\n        port => 5044\n    }\n}\nfilter {\n    if \"eu.gcr.io/sk-private-registry/skysoft-atm/\" not in [kubernetes][container][image] and \"elasticsearch\" not in [kubernetes][labels][name] {\n        drop {}\n    }\n    if [kubernetes][labels][name] == \"elasticsearch\" {\n        grok {\n            match => [ \"message\", \"[%{TIMESTAMP_ISO8601:timestamp}][%{DATA:level}%{SPACE}][%{DATA:source}%{SPACE}]%{SPACE}%{GREEDYDATA:message}\" ]\n            overwrite => [ \"message\" ]\n        }\n    }\n    if [kubernetes][labels][name] == \"ems\" {\n        grok {\n            match => { \"message\" => \"%{TIMESTAMP_ISO8601:logdate} %{LOGLEVEL:level}: %{GREEDYDATA:message}\" }\n            overwrite => [\"message\"]\n        \n        date {\n            match => [ \"logdate\", \"ISO8601\", \"yyyy-MM-dd HH:mm:ss,SSS\", \"yyyy-MM-dd HH:mm:ss.SSS\" ]\n            remove_field => [ \"logdate\" ]\n        }\n    }\n    date {\n        match => [ \"timestamp\", \"ISO8601\", \"yyyy-MM-dd HH:mm:ss,SSS\", \"yyyy-MM-dd HH:mm:ss.SSS\" ]\n    }\n}\noutput {\n    elasticsearch {\n        index => \"filebeat-%{+yyyy.MM.dd}\"\n        cloud_id => \"Test\"\n        cloud_auth => \"Test\"\n    }\n}",
		"settings": {
			"pipeline.batch.delay": 50,
			"pipeline.batch.size": 125,
			"pipeline.workers": 1,
			"queue.checkpoint.writes": 1024,
			"queue.max_bytes": "1gb",
			"queue.type": "persisted"
		}
	}`)
	c           *Client
	pipelineRef LogstashConfiguration
)

func init() {
	// Check env variables exist
	assertPrerequisites()

	c = NewClient(os.Getenv("CLOUD_AUTH"), os.Getenv("KIBANA_URL"))
	err := json.Unmarshal(refPipeline, &pipelineRef)
	if err != nil {
		panic("error trying to load pipeline definition")
	}
}

func TestCreateAndGetPipeline(t *testing.T) {
	ctx := context.Background()
	pipeline := &LogstashPipeline{ID: generatePipelineID(), Configuration: &pipelineRef}

	err := c.CreateOrUpdateLogstashPipeline(ctx, pipeline)
	assert.Nil(t, err, "[ Creation ] expecting nil error")

	res, err := c.GetLogstashPipeline(ctx, pipeline.ID)
	assert.Nil(t, err, "[ Reading ] expecting nil error")
	t.Logf("Here is the settings object %v", res.Configuration.Settings)
	assert.Equal(t, pipeline.ID, res.ID, "expecting same IDs")
	assert.Equal(t, pipeline.Configuration.Description, res.Configuration.Description, "expecting same description")
	assert.Equal(t, pipeline.Configuration.Settings, res.Configuration.Settings, "expecting same settings")
	assert.Equal(t, pipeline.Configuration.Pipeline, res.Configuration.Pipeline, "expecting same pipeline definition")
	// username is generated based on token used.
	assert.Equal(t, "elastic", res.Configuration.Username, "expecting same username")

	// Clean environment
	err = c.DeleteLogstashPipeline(ctx, pipeline.ID)
	assert.Nil(t, err, "[ Deleting ] expecting nil error")
}

func TestEmptySettings(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		p             *LogstashPipeline
		expectedError error
	}{
		{NewLogstashPipeline("", "Description", "Test", nil), fmt.Errorf("ID cannot be empty")},
		{NewLogstashPipeline(generatePipelineID(), "", "", nil), fmt.Errorf("pipeline definition cannot be empty")},
	}

	for _, test := range tests {
		err := c.CreateOrUpdateLogstashPipeline(ctx, test.p)
		if assert.Error(t, err) {
			assert.Equal(t, test.expectedError, err)
		}
	}
}

func TestCleanURL(t *testing.T) {
	tests := []struct {
		baseURL  string
		suffix   string
		expected string
	}{
		{"https://myurl.com/", "/directive", "https://myurl.com/directive"},
		{"https://myurl.com", "directive", "https://myurl.com/directive"},
		{"https://myurl.com/", "directive", "https://myurl.com/directive"},
		{"https://myurl.com", "/directive", "https://myurl.com/directive"},
	}

	for _, test := range tests {
		result := cleanURL(test.baseURL, test.suffix)
		assert.Equal(t, test.expected, result, "expected (%s) and result (%s) should be equal", test.expected, result)
	}
}

func TestGetAll(t *testing.T) {
	pipelines, err := c.GetLogstashPipelines(context.Background())
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, pipelines.Pipelines, "pipelines slice should not be empty")
	// At this stage no expectation on the number of pipelines stored
	for _, pipe := range pipelines.Pipelines {
		assert.NotEmpty(t, pipe.ID)
		assert.NotEmpty(t, pipe.LastModified)
		assert.NotEmpty(t, pipe.Username)
	}
}

func assertPrerequisites() {
	if len(os.Getenv("CLOUD_AUTH")) == 0 {
		panic("CLOUD_AUTH env variable should be set for tests")
	}
	if len(os.Getenv("KIBANA_URL")) == 0 {
		panic("KIBANA_URL env variable should be set for tests")
	}
}

func generatePipelineID() string {
	return shortuuid.New()
}
