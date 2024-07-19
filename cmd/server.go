package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rcbadiale/go-rate-limiter/internal/stores"
	"github.com/rcbadiale/go-rate-limiter/pkg/config"
	"github.com/rcbadiale/go-rate-limiter/pkg/limiter"
	"github.com/rcbadiale/go-rate-limiter/pkg/middlewares"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, will use environment variables")
	}
	cfg := config.LoadConfig()
	store := stores.NewMemory()
	ipLimiter := limiter.NewLimiter(store,
		cfg.IPLimit,
		cfg.IPDuration,
	)
	apiKeyLimiter := limiter.NewLimiter(store,
		cfg.APIKeyLimit,
		cfg.APIKeyDuration,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloRoute)

	log.Println("Server started on port 8080")
	mid := middlewares.LogRequest(mux)
	ipRateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(ipLimiter, nil)
	apiKeyRateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(apiKeyLimiter, keyMapper)
	mid = ipRateLimiterMiddleware(mid)
	mid = apiKeyRateLimiterMiddleware(mid)
	http.ListenAndServe(":8080", mid)
	log.Println("Server stopped")
}

func helloRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Bem-vindo à API com Rate Limiter!"}`))
}

func keyMapper(r *http.Request) string {
	api_key := r.Header.Get("API_KEY")
	if api_key != "" {
		return fmt.Sprintf("API_KEY:%v", api_key)
	}
	return ""
}
