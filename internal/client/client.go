package client

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v7"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

type Client struct {
	Config *Agent
	*http.Client
}

type Agent struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PoolInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

var (
	address        = flag.String("a", "localhost:8080", "server address")
	reportInterval = flag.Duration("r", 10*time.Second, "report interval")
	poolInterval   = flag.Duration("p", 2*time.Second, "pool interval")
)

func InitConfigAgent() *Agent {
	cfg := &Agent{}

	flag.Parse()
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address = os.Getenv("ADDRESS"); cfg.Address == "" {
		cfg.Address = *address
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval == "" {
		if cfg.ReportInterval, err = time.ParseDuration(envReportInterval); err != nil {
			cfg.ReportInterval = *reportInterval
		}
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval == "" {
		if cfg.PoolInterval, err = time.ParseDuration(envPoolInterval); err != nil {
			cfg.PoolInterval = *poolInterval
		}
	}

	return cfg
}

func NewHTTPClient(cfg *Agent) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		Config: cfg,
		Client: &http.Client{
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
		Host:   c.Config.Address,
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
		req.Header.Set("Accept-Encoding", "gzip")
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
