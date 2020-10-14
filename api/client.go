package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/skysoft-atm/terraform-provider-elastic/utils"
)

// Client is the high-level structure to interact with Elastic API
type Client struct {
	BaseURL    string
	cloudAuth  string
	HTTPClient *http.Client
}

type errorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

// LogstashPipelines object retrieved via the /pipelines directive
type LogstashPipelines struct {
	Pipelines []struct {
		ID           string `json:"id"`
		Description  string `json:"description,omitempty"`
		LastModified string `json:"last_modified,omitempty"`
		Username     string `json:"username"`
	} `json:"pipelines"`
}

// LogstashPipeline object to be used with elastic API to define logstash pipelines
// https://www.elastic.co/guide/en/kibana/current/logstash-configuration-management-api-create.html#logstash-configuration-management-api-create-request-body
type LogstashPipeline struct {
	ID            string                 `json:"id"`
	Configuration *LogstashConfiguration `json:"config,omitempty"`
}

// LogstashConfiguration is the underlying struct sent via kibana API (ID should not be included)
type LogstashConfiguration struct {
	Description string    `json:"description,omitempty"`
	Username    string    `json:"username,omitempty"`
	Pipeline    string    `json:"pipeline"`
	Settings    *Settings `json:"settings,omitempty"`
}

// Settings defines the options for the logstash pipeline
type Settings struct {
	PipelineBatchDelay    int    `json:"pipeline.batch.delay,omitempty"`
	PipelineBatchSize     int    `json:"pipeline.batch.size,omitempty"`
	PipelineWorkers       int    `json:"pipeline.workers,omitempty"`
	QueueCheckpointWrites int    `json:"queue.checkpoint.writes,omitempty"`
	QueueMaxBytes         string `json:"queue.max_bytes,omitempty"`
	QueueType             string `json:"queue.type,omitempty"`
}

// NewLogstashPipeline returns a *LogstashPipeline struct
func NewLogstashPipeline(id, description, pipeline string, settings *Settings) *LogstashPipeline {
	return &LogstashPipeline{
		ID: id,
		Configuration: &LogstashConfiguration{
			Description: description,
			Pipeline:    pipeline,
			Settings:    settings,
		},
	}
}

// NewLogstashPipelineSettings returns a *Settings struct
func NewLogstashPipelineSettings(batchDelay, batchSize, workers, queueCheckpointWrites int, queueMaxBytes, queueType string) *Settings {
	return &Settings{
		PipelineBatchDelay:    batchDelay,
		PipelineBatchSize:     batchSize,
		PipelineWorkers:       workers,
		QueueCheckpointWrites: queueCheckpointWrites,
		QueueMaxBytes:         queueMaxBytes,
		QueueType:             queueType,
	}
}

// NewClient returns a new HTTP Client
func NewClient(cloudAuth string, kibanaURL string) *Client {
	return &Client{
		BaseURL:   kibanaURL,
		cloudAuth: cloudAuth,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

const (
	crudBaseURL   = "/api/logstash/pipeline"
	getAllBaseURL = "/api/logstash/pipelines"
)

// GetLogstashPipelines return the current list of pipelines
func (c *Client) GetLogstashPipelines(ctx context.Context) (*LogstashPipelines, error) {
	url := cleanURL(c.BaseURL, getAllBaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := LogstashPipelines{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetLogstashPipeline retrieve the pipeline identified with the unique ID
func (c *Client) GetLogstashPipeline(ctx context.Context, id string) (*LogstashPipeline, error) {
	url := cleanURL(cleanURL(c.BaseURL, crudBaseURL), id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		if err.Error() == "Not Found" {
			return &LogstashPipeline{}, nil
		}
		return nil, err
	}

	req = req.WithContext(ctx)

	res := LogstashConfiguration{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &LogstashPipeline{
		ID:            id,
		Configuration: &res,
	}, nil
}

// DeleteLogstashPipeline deletes a specific logstash pipeline
func (c *Client) DeleteLogstashPipeline(ctx context.Context, id string) error {
	url := cleanURL(cleanURL(c.BaseURL, crudBaseURL), id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// CreateOrUpdateLogstashPipeline creates/updates a specific logstash pipeline
func (c *Client) CreateOrUpdateLogstashPipeline(ctx context.Context, lp *LogstashPipeline) error {

	err := checkPrerequisites(lp)
	if err != nil {
		return err
	}

	url := cleanURL(cleanURL(c.BaseURL, crudBaseURL), lp.ID)

	// marshal LogstashPipeline to JSON
	json, err := json.Marshal(lp.Configuration)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("kbn-xsrf", "true")
	req.Header.Set("Cache-Control", "no-cache")

	username, password, err := utils.ParseTwoPartID(c.cloudAuth, "username", "password")
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		errBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}
		return fmt.Errorf("unknown error, status code: %d, message: %s", res.StatusCode, string(errBody))
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
	}

	return nil
}

func (p *LogstashPipeline) String() string {
	js, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return ""
	}
	return string(js)
}

func cleanURL(baseURL, suffix string) string {
	if strings.HasSuffix(baseURL, "/") && strings.HasPrefix(suffix, "/") {
		return fmt.Sprintf("%s%s", baseURL, trimFirstRune(suffix))
	}
	if !strings.HasSuffix(baseURL, "/") && !strings.HasPrefix(suffix, "/") {
		return fmt.Sprintf("%s/%s", baseURL, suffix)
	}
	return fmt.Sprintf("%s%s", baseURL, suffix)
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func checkPrerequisites(p *LogstashPipeline) error {
	if len(p.ID) == 0 {
		return fmt.Errorf("ID cannot be empty")
	}
	if len(p.Configuration.Pipeline) == 0 {
		return fmt.Errorf("pipeline definition cannot be empty")
	}
	return nil
}
