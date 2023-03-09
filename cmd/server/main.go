package main

import (
	"github.com/eugeniylennik/alertics/internal/router"
	"github.com/eugeniylennik/alertics/internal/server"
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()
	s := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error %v", err)
		}
	}()
	server.ShutdownServer(s)
}
