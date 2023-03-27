package client

import (
	"bytes"
	"encoding/json"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	host = "127.0.0.1"
	port = "8080"
)

type Client struct {
	*http.Client
}

func NewHTTPClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		&http.Client{
			Timeout: time.Second * 5,
			Transport: &http.Transport{
				MaxIdleConns: 20,
			},
			Jar: jar,
		},
	}, nil
}

func (c *Client) SendMetrics(m metrics.ListMetrics) error {
	if len(m) == 0 {
		return nil
	}
	addr := url.URL{
		Scheme: "http",
		Host:   host + ":" + port,
		Path:   "/update/",
	}
	for _, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		req, err := http.NewRequest(
			http.MethodPost,
			addr.String(),
			bytes.NewBuffer(b),
		)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("error closing response body: %v", err)
			}
		}()
	}

	return nil
}
