package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8081", "server port")
	flag.Parse()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from backend on port %s\n", *port)
	})

	log.Printf("Backend running on :%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, handler))
}
