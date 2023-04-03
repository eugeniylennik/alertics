package main

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/router"
	"github.com/eugeniylennik/alertics/internal/server"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/eugeniylennik/alertics/internal/storage/file"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	store := storage.NewMemStorage()
	r := router.NewRouter(store)

	s := &http.Server{
		Addr:    server.Config.Address,
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
			log.Printf("HTTP server gracefully stopped\n")
		}
	}
}

func collectMetricsToFile(ctx context.Context, store *storage.MemStorage) error {
	interval := time.NewTicker(server.Config.StoreInterval)
	defer interval.Stop()

	w, err := file.NewWriter(server.Config.StoreFile)
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
	if server.Config.Restore {
		r, err := file.NewReader(server.Config.StoreFile)
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
