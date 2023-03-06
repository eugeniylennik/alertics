package main

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/client"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"log"
	"time"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tPool := time.NewTicker(pollInterval)
	tReport := time.NewTicker(reportInterval)

	c, err := client.NewHttpClient()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan []metrics.Data, 0)
	go func() {
		for range tPool.C {
			ch <- metrics.CollectMetrics()
		}
	}()

	for range tReport.C {
		if err := c.SendMetrics(ctx, <-ch); err != nil {
			log.Fatal(err)
		}
	}
}
