package main

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := router.NewRouter()

	s := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		)
		<-sigint

		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error %v", err)
	}

	<-idleConnsClosed
}
