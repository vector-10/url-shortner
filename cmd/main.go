package main

import (
	"fmt"
	"net/http"
	"syscall"
	"time"
	"context"
	"os"
	"os/signal"	
	"github.com/vector-10/url-shortner/internal/handler"
	"github.com/vector-10/url-shortner/internal/store"
)

func main() {
	s := store.NewMemoryStore()
	s.StartCleanup(1 * time.Minute)
	h := handler.NewHandler(s)


	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.ShortenURL)
	mux.HandleFunc("GET /{slug}/qr", h.QRCode)
	mux.HandleFunc("GET /{slug}", h.Redirect)
	mux.HandleFunc("GET /{slug}/stats", h.Stats)

	server := &http.Server {
		Addr: ":8080",
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server running on PORT 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Forced shutdown:", err)
	}

	fmt.Println("Server stopped")

}