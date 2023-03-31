package main

import (
	"context"
	"github.com/caarlos0/env/v7"
	"github.com/eugeniylennik/alertics/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
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
	r := router.NewRouter()

	s := &http.Server{
		Addr:    Config.Address,
		Handler: r,
	}

	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
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
