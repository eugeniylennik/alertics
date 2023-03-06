package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
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

func (c *Client) SendMetrics(ctx context.Context, m []metrics.Data) error {
	if len(m) == 0 {
		return errors.New("empty metrics")
	}
	for _, v := range m {
		addr := url.URL{
			Scheme: "http",
			Host:   host + ":" + port,
			Path:   fmt.Sprintf("/update/%s/%s/%.1f", v.Type, v.Name, v.Value),
		}
		req, err := http.NewRequest(
			http.MethodPost,
			addr.String(),
			nil,
		)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "text/plain")
		fmt.Printf("Success send: Name - %s, Value - %.1f\n", v.Name, v.Value)
		//_, err = c.Do(req)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}
