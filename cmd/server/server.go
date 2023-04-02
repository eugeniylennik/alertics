package main

import (
	"context"
	"github.com/caarlos0/env/v7"
	"github.com/eugeniylennik/alertics/internal/router"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/eugeniylennik/alertics/internal/storage/file"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"5s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

var Config Server

func init() {
	Config = Server{}
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	store := storage.NewMemStorage()
	r := router.NewRouter(store)

	s := &http.Server{
		Addr:    Config.Address,
		Handler: r,
	}

	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := restoreMetrics(store); err != nil {
			log.Println(err)
		}
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	go func() {
		if err := collectMetricsToFile(ctx, store); err != nil {
			errChan <- err
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("HTTP server ListenAndServe Error %v", err)
		cancel()
		return
	case <-sig:
		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v\n", err)
			os.Exit(1)
		} else {
			//collectMetricsToFile(ctx, store)
			log.Printf("HTTP server gracefully stopped\n")
		}
	}
}

func collectMetricsToFile(ctx context.Context, store *storage.MemStorage) error {
	interval := time.NewTicker(Config.StoreInterval)
	defer interval.Stop()

	w, err := file.NewWriter(Config.StoreFile)
	if err != nil {
		return err
	}

	for {
		select {
		case <-interval.C:
			mBz, err := store.GetAllMetrics()
			if err != nil {
				return err
			}
			if err := w.WriteMetrics(mBz); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func restoreMetrics(store *storage.MemStorage) error {
	if Config.Restore {
		r, err := file.NewReader(Config.StoreFile)
		if err != nil {
			return err
		}
		m, err := r.ReadMetrics()
		if err != nil {
			return err
		}
		for _, v := range m {
			switch v.Type {
			case storage.Gauge:
				_ = store.AddGauge(v)
			case storage.Counter:
				_ = store.AddCounter(v)
			}
		}
	}
	return nil
}
