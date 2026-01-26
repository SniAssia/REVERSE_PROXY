package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"REVERSE_PROXY/cmd/backend"
	"REVERSE_PROXY/internal/admin"
	"REVERSE_PROXY/internal/health"
	"REVERSE_PROXY/internal/loadbalancer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := backend.NewServer() 

	var lbValue atomic.Value
	lbValue.Store(loadbalancer.Newround(server)) // default strategy

	checker := health.NewChecker(server, 5*time.Second)
	go checker.Start(ctx)

	proxyHandler := NewProxyHandler(server, &lbValue)
	proxyServer := &http.Server{
		Addr:    ":8080",
		Handler: proxyHandler,
	}

	go func() {
		log.Println("[proxy] listening on :8080")
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[proxy] server error: %v", err)
		}
	}()

	adminServer := admin.Newserver(":9000", server, &lbValue)
	go func() {
		log.Println("[admin] listening on :9000")
		if err := adminServer.Start(); err != nil {
			log.Fatalf("[admin] server error: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("[main] shutdown signal received")
	cancel() 

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := proxyServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("[proxy] shutdown error: %v", err)
	}
	if err := adminServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("[admin] shutdown error: %v", err)
	}

	log.Println("[main] shutdown complete")
}
