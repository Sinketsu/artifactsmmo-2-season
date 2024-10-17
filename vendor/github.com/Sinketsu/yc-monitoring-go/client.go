package ycmonitoringgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	url       string
	queryArgs string
	token     string

	client *http.Client
	logger Logger
}

func NewClient(folder string, token string, options ...Option) *Client {
	opts := defaultOptions()

	for _, apply := range options {
		apply(opts)
	}

	queryArgs := url.Values{}
	queryArgs.Add("folderId", folder)
	queryArgs.Add("service", "custom")

	return &Client{
		url:       opts.url,
		queryArgs: queryArgs.Encode(),
		token:     token,

		client: opts.client,
		logger: opts.logger,
	}
}

func (c *Client) Run(ctx context.Context, registry *Registry, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.send(ctx, registry)
			if err != nil && c.logger != nil {
				c.logger.Error("fail send metrics", err)
			}
		}
	}
}

func (c *Client) send(ctx context.Context, registry *Registry) error {
	r := Request{}

	registry.Range(func(i int, m Metric) {
		r.Metrics = append(r.Metrics, m.GetMetrics()...)
	})

	body, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("fail marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/v2/data/write", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("fail to create new http request: %w", err)
	}

	req.URL.RawQuery = c.queryArgs
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("fail to do http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("got non 200 code: %d", resp.StatusCode)
	}

	return nil
}
