package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ShutdownServer(s *http.Server) {
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
