package main

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/client"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c, err := client.NewHTTPClient()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan []metrics.Data)

	go collectMetrics(ctx, ch)
	go sendMetrics(ctx, c, ch)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-s:
		cancel()
	case <-ctx.Done():
	}
}

func collectMetrics(ctx context.Context, ch chan []metrics.Data) {
	tPool := time.NewTicker(client.Config.PoolInterval)
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
	tReport := time.NewTicker(client.Config.ReportInterval)
	defer tReport.Stop()

	var m []metrics.Data
	for {
		select {
		case newM := <-ch:
			m = newM
		case <-tReport.C:
			if err := c.SendMetrics(m); err != nil {
				log.Fatal(err)
			}
			m = nil
		case <-ctx.Done():
			return
		}
	}
}
