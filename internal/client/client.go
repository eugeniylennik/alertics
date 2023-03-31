package client

import (
	"bytes"
	"encoding/json"
	"github.com/caarlos0/env/v7"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	*http.Client
}

type Agent struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PoolInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

var Config Agent

func init() {
	Config = Agent{}
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}
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

func (c *Client) SendMetrics(d []metrics.Data) error {
	if len(d) == 0 {
		return nil
	}
	addr := url.URL{
		Scheme: "http",
		Host:   Config.Address,
		Path:   "/update",
	}
	for _, v := range d {
		m := metrics.Metrics{
			ID:    v.Name,
			MType: v.Type,
		}

		if v.Type == storage.Gauge {
			m.Value = &v.Value
		} else {
			i := int64(v.Value)
			m.Delta = &i
		}

		b, err := json.Marshal(m)
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
