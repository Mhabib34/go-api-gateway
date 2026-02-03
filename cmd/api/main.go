package main

import (
	"fmt"
	"log"
	"net/http"

	"go_api_gateway/config"
	"go_api_gateway/internal/middleware"
	"go_api_gateway/internal/proxy"
)

func main() {

	// =============================
	// Backend services
	// =============================
	backends := []*proxy.Backend{
		{URL: config.MustParseURL("http://localhost:8081"), Alive: true},
		{URL: config.MustParseURL("http://localhost:8082"), Alive: true},
	}

	// =============================
	// Load Balancer
	// =============================
	loadBalancer := proxy.NewLoadBalancer(backends)

	// =============================
	// Redis (Rate Limiter)
	// =============================
	redisClient := config.NewRedisClient()

	// =============================
	// Middleware Chain
	// =============================
	handler := middleware.RateLimiterByIP(redisClient)(loadBalancer)

	// =============================
	// Run Server
	// =============================
	fmt.Println("🚀 API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
