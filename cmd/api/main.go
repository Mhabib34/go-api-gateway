package main

import (
	"go_api_gateway/config"
	"go_api_gateway/internal/proxy"
	"log"
	"net/http"
	"time"
)

func main() {
	backends := []*proxy.Backend{
		{URL: config.MustParseURL("http://localhost:8081"), Alive: true},
		{URL: config.MustParseURL("http://localhost:8082"), Alive: true},
	}

	client := proxy.NewClient("payment-service")

	for range 10 {
		_, err := client.Get("https://httpstat.us/500")
		log.Println("err:", err)
		time.Sleep(500 * time.Millisecond)
	}


	lb := proxy.NewLoadBalancer(backends)

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", lb))
}
