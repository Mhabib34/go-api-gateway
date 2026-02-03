package main

import (
	"log"
	"time"

	"go_api_gateway/internal/proxy"
)

func main() {

	client := proxy.NewClient("testing-service")

	for range 10 {
		_, err := client.Get("https://httpstat.us/500")
		log.Println("err:", err)
		time.Sleep(500 * time.Millisecond)
	}
}
