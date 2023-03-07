package main

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := &http.Server{
		Addr: "localhost:8080",
	}

	h := handlers.NewStorage()
	http.HandleFunc("/update/", h.RecordMetrics)

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error %v", err)
		}
	}()

	shutdownServer(s)
}

func shutdownServer(s *http.Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	<-signalChan
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
		os.Exit(1)
	} else {
		log.Printf("HTTP server gracefully stopped\n")
	}
}
