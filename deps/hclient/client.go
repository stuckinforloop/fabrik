package hclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(timeout int64) *Client {
	if timeout == 0 {
		timeout = 5
	}

	return &Client{
		&http.Client{
			Timeout: time.Second * time.Duration(timeout),
		},
	}
}

func (c *Client) Do(ctx context.Context, url, method string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}
