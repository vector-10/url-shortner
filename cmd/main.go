package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vector-10/url-shortner/internal/handler"
	"github.com/vector-10/url-shortner/internal/store"
)


func corsMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping postgres:", err)
	}

	s := store.NewPostgresStore(db)
	us := store.NewPostgresUserStore(db)
	cache := store.NewRedisCache(os.Getenv("REDIS_ADDR"))

	h := handler.NewHandler(s, cache)
	ah := handler.NewAuthHandler(us)
	oh := handler.NewOAuthHandler(us)
 

	mux := http.NewServeMux()
	
	//public routes
	mux.HandleFunc("POST /signup", ah.Signup)
	mux.HandleFunc("POST /login", ah.Login)
	mux.HandleFunc("GET /{slug}", h.Redirect)
	mux.HandleFunc("GET /{slug}/qr", (h.QRCode))
	mux.HandleFunc("GET /auth/google", oh.GoogleLogin)
	mux.HandleFunc("GET /auth/google/callback", oh.GoogleCallback)

	//protected routes
	mux.HandleFunc("POST /shorten", handler.RequireAuth(h.ShortenURL))

	mux.HandleFunc("GET /urls", handler.RequireAuth(h.ListURLs))	
	mux.HandleFunc("GET /{slug}/stats", handler.RequireAuth(h.Stats))

	server := &http.Server {
		Addr: ":8080",
		Handler: corsMiddleware(mux),
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