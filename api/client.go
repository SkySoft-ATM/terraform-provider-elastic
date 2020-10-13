package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/copier"

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

// CreateOrUpdateLogstashPipeline creates/updates a specific logstash pipeline
func (c *Client) CreateOrUpdateLogstashPipeline(ctx context.Context, logstashPipeline *LogstashPipeline, ID string) error {
	// Trick to avoid "definition for this key is missing exception"
	pipeline := LogstashPipeline{}

	if len(logstashPipeline.ID) > 0 {
		copier.Copy(&pipeline, logstashPipeline)
		pipeline.ID = ""
	}

	// marshal LogstashPipeline to json
	json, err := json.Marshal(&pipeline)
	log.Printf("Voici le message en JSON %s", json)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", c.BaseURL, ID), bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	log.Printf("Voici la requete %v", req)

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

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
	}

	return nil
}
