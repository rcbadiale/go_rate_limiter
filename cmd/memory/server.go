package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rcbadiale/go-rate-limiter/internal/stores/memory"
	"github.com/rcbadiale/go-rate-limiter/pkg/config"
	"github.com/rcbadiale/go-rate-limiter/pkg/limiter"
	"github.com/rcbadiale/go-rate-limiter/pkg/middlewares"
)

func main() {
	cfg := config.LoadConfig()
	store := memory.NewMemoryStore()

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
	ipRateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(ipLimiter, nil)
	apiKeyRateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(apiKeyLimiter, keyMapper)
	mid := ipRateLimiterMiddleware(mux)
	mid = apiKeyRateLimiterMiddleware(mid)
	mid = middlewares.LogRequest(mid)
	http.ListenAndServe(":8080", mid)
	log.Println("Server stopped")
}

func helloRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to the Rate Limiter API!"}`))
}

func keyMapper(r *http.Request) string {
	apiKey := r.Header.Get("API_KEY")
	if apiKey != "" {
		return fmt.Sprintf("API_KEY:%v", apiKey)
	}
	return ""
}
