package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client is the high-level structure to interact with Elastic API
type Client struct {
	BaseURL    string
	cloudAuth  string
	HTTPClient *http.Client
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	StatusCode int         `json:"code"`
	Data       interface{} `json:"data"`
}

// LogstashPipeline object to be used with elastic API to define logstash pipelines
// https://www.elastic.co/guide/en/kibana/current/logstash-configuration-management-api-create.html#logstash-configuration-management-api-create-request-body
type LogstashPipeline struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	Username    string `json:"username,omitempty"`
	Pipeline    string `json:"pipeline,omitempty"`
	Settings    struct {
		PipelineBatchDelay    int    `json:"pipeline.batch.delay,omitempty"`
		PipelineBatchSize     int    `json:"pipeline.batch.size,omitempty"`
		PipelineWorkers       int    `json:"pipeline.workers,omitempty"`
		QueueCheckpointWrites int    `json:"queue.checkpoint.writes,omitempty"`
		QueueMaxBytes         string `json:"queue.max_bytes,omitempty"`
		QueueType             string `json:"queue.type,omitempty"`
	} `json:"settings,omitempty"`
}

// NewClient instantiates a new HTTP client
func NewClient(cloudAuth string, kibanaURL string) *Client {
	return &Client{
		BaseURL:   fmt.Sprintf("%s/api/logstash/pipeline", kibanaURL),
		cloudAuth: cloudAuth,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetLogstashPipeline retrieve the pipeline identified with the unique ID
func (c *Client) GetLogstashPipeline(ctx context.Context, ID string) (*LogstashPipeline, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.BaseURL, ID), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := LogstashPipeline{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateLogstashPipeline updates a specific logstash pipeline
func (c *Client) UpdateLogstashPipeline(ctx context.Context) error {
	return nil
}

// DeleteLogstashPipeline deletes a specific logstash pipeline
func (c *Client) DeleteLogstashPipeline(ctx context.Context) error {
	return nil
}

// CreateLogstashPipeline creates a specific logstash pipeline
func (c *Client) CreateLogstashPipeline(ctx context.Context) error {
	return nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("kbn-xsrf", "true")
	req.Header.Set("Cache-Control", "no-cache")
	auth := strings.Split(c.cloudAuth, ":")
	req.SetBasicAuth(auth[0], auth[1])

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

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
