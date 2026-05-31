package main

import (
	"fmt"
	"net/http"
	"syscall"
	"time"
	"context"
	"os"
	"os/signal"	

	"github.com/joho/godotenv"
	"github.com/vector-10/url-shortner/internal/handler"
	"github.com/vector-10/url-shortner/internal/store"
)

func main() {
	godotenv.Load()

	// s := store.NewMemoryStore()
	// s.StartCleanup(1 * time.Minute)
	redisAddr := os.Getenv("REDIS_ADDR")


	s := store.NewRedisStore(redisAddr)
	us := store.NewRedisUserStore(redisAddr)

	h := handler.NewHandler(s)
	ah := handler.NewAuthHandler(us)
	oh := handler.NewOAuthHandler(us)


	mux := http.NewServeMux()
	
	//public routes
	mux.HandleFunc("POST /signup", ah.Signup)
	mux.HandleFunc("POST/login", ah.Login)
	mux.HandleFunc("GET /{slug}", h.Redirect)
	mux.HandleFunc("GET /auth/google", oh.GoogleLogin)
	mux.HandleFunc("GET /auth/google/callback", oh.GoogleCallback)

	//protected routes
	mux.HandleFunc("POST /shorten", handler.RequireAuth(h.ShortenURL))
	mux.HandleFunc("GET /{slug}/qr", handler.RequireAuth(h.QRCode))
	 mux.HandleFunc("GET /urls", handler.RequireAuth(h.ListURLs))	
	mux.HandleFunc("GET /{slug}/stats", handler.RequireAuth(h.Stats))

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