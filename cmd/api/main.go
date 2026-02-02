package main

import (
	"go_api_gateway/config"
	"go_api_gateway/internal/proxy"
	"log"
	"net/http"
)

func main() {
	backends := []*proxy.Backend{
		{URL: config.MustParseURL("http://localhost:8081"), Alive: true},
		{URL: config.MustParseURL("http://localhost:8082"), Alive: true},
	}

	lb := proxy.NewLoadBalancer(backends)

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", lb))
}
