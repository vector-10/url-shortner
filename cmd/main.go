package main

import (
	"fmt"
	"net/http"
	"time"
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

	fmt.Println("Server running on PORT 8080")
	http.ListenAndServe(":8080", mux)
}