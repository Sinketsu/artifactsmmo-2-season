package monitoring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	globalMetrics []Metric
	globalMu      sync.Mutex
)

type Request struct {
	Metrics []metric `json:"metrics"`
}

type Metric interface {
	Get() []metric
}

type Client struct {
	url    string
	token  string
	folder string

	cli *http.Client
}

func NewClient(url string, folder string, token string) *Client {
	return &Client{
		url:    url,
		token:  token,
		folder: folder,

		cli: http.DefaultClient,
	}
}

func (c *Client) send() error {
	globalMu.Lock()
	defer globalMu.Unlock()

	r := Request{}
	for _, m := range globalMetrics {
		r.Metrics = append(r.Metrics, m.Get()...)
	}

	body, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("fail marshal json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("fail to create new http request: %w", err)
	}

	q := url.Values{}
	q.Add("folderId", c.folder)
	q.Add("service", "custom")

	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cli.Do(req)
	if err != nil {
		return fmt.Errorf("fail to do http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("non 200 code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Run(interval time.Duration) {
	for range time.Tick(interval) {
		err := c.send()
		if err != nil {
			fmt.Println("fail to send metrics:", err.Error())
		}
	}
}
