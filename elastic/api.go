package elastic

import (
	"bytes"
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
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

// LogstashPipeline object to be used with elastic API to define logstash pipelines
// https://www.elastic.co/guide/en/kibana/current/logstash-configuration-management-api-create.html#logstash-configuration-management-api-create-request-body
type LogstashPipeline struct {
	ID          string    `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	Username    string    `json:"username,omitempty"`
	Pipeline    string    `json:"pipeline,omitempty"`
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
		ID:          id,
		Description: description,
		Pipeline:    pipeline,
		Settings:    settings,
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
		BaseURL:   fmt.Sprintf("%s/api/logstash/pipeline", kibanaURL),
		cloudAuth: cloudAuth,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetLogstashPipeline retrieve the pipeline identified with the unique ID
func (c *Client) GetLogstashPipeline(ctx context.Context, ID string) (*LogstashPipeline, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.BaseURL, ID), nil)
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
func (c *Client) UpdateLogstashPipeline(ctx context.Context, logstashPipeline *LogstashPipeline, ID string) error {
	// marshal LogstashPipeline to json
	json, err := json.Marshal(logstashPipeline)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", c.BaseURL, ID), bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// DeleteLogstashPipeline deletes a specific logstash pipeline
func (c *Client) DeleteLogstashPipeline(ctx context.Context, ID string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.BaseURL, ID), nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// CreateLogstashPipeline creates a specific logstash pipeline
func (c *Client) CreateLogstashPipeline(ctx context.Context, logstashPipeline *LogstashPipeline, ID string) error {
	// marshal LogstashPipeline to json
	json, err := json.Marshal(logstashPipeline)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", c.BaseURL, ID), bytes.NewBuffer(json))
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
	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
	}

	return nil
}
