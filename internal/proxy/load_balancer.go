package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	mu    sync.RWMutex
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

type LoadBalancer struct {
	backends []*Backend
	current  int
	mu       sync.Mutex
}

func NewLoadBalancer(backends []*Backend) *LoadBalancer {
	return &LoadBalancer{
		backends: backends,
	}
} 

func (lb *LoadBalancer) nextBackend() *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := 0; i < len(lb.backends); i++ {
		idx := (lb.current + i) % len(lb.backends)
		if lb.backends[idx].IsAlive() {
			lb.current = (idx + 1) % len(lb.backends)
			return lb.backends[idx]
		}
	}
	return nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.nextBackend()
	if backend == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		backend.SetAlive(false)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}
