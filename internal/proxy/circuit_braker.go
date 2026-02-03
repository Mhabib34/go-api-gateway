package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker/v2"
)

// Client
type Client struct {
	cb *gobreaker.CircuitBreaker[[]byte]
	client *http.Client
}

func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker[[]byte] {
	// Circuit Breaker
	st := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Timeout:     10 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},

		// OnStateChange
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("[CB] %s: %s → %s\n", name, from, to)
		},
	}

	return gobreaker.NewCircuitBreaker[[]byte](st)
}


func NewClient(name string) *Client {
	// Client
	return &Client{
		cb: NewCircuitBreaker(name),
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}


func (c *Client) Get(url string) ([]byte, error) {
	// Circuit Breaker
	return c.cb.Execute(func() ([]byte, error) {
		// Request
		resp, err := c.client.Get(url)
		if err != nil {
			return nil, err
		}
		
		// Response
		defer resp.Body.Close()
		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		// Body
		return io.ReadAll(resp.Body)
	})
}

