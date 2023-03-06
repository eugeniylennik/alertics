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

	c, err := client.NewHttpClient()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan []metrics.Data)

	go collectMetrics(ctx, ch)
	go sendMetrics(ctx, c, ch)

	<-ctx.Done()
}

func collectMetrics(ctx context.Context, ch chan []metrics.Data) {
	tPool := time.NewTicker(pollInterval)
	defer tPool.Stop()

	for {
		select {
		case <-tPool.C:
			ch <- metrics.CollectMetrics()
		case <-ctx.Done():
			return
		}
	}
}

func sendMetrics(ctx context.Context, c *client.Client, ch chan []metrics.Data) {
	tReport := time.NewTicker(reportInterval)
	defer tReport.Stop()

	var m []metrics.Data
	for {
		select {
		case newM := <-ch:
			m = newM
		case <-tReport.C:
			if len(m) > 0 {
				if err := c.SendMetrics(ctx, m); err != nil {
					log.Fatal(err)
				}
				m = nil
			}
		case <-ctx.Done():
			return
		}
	}
}
