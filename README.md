# 🚀 Go API Gateway

A lightweight, production-ready **API Gateway** built in Go. This project provides a single entry point for backend services with built-in **load balancing**, **rate limiting**, and **circuit breaker** capabilities — no heavy frameworks, just clean Go and proven libraries.

---

## ✨ Features

- **Round-Robin Load Balancer** — distributes traffic across multiple backend services with automatic health-based skipping.
- **Rate Limiter** — IP-based rate limiting using Redis as a shared state store (default: 100 requests/minute per IP).
- **Circuit Breaker** — per-backend circuit breaker powered by `sony/gobreaker`, trips after a 60% failure ratio on 5+ requests, and auto-recovers after 10 seconds.
- **Reverse Proxy** — uses Go's built-in `httputil.ReverseProxy` for low-overhead proxying to backends.

---

## 📁 Project Structure

```
go-api-gateway/
├── cmd/
│   └── backend/  
│       └── main.go                # Gateway entry point
├── config/
│   ├── backend.go                 # URL parsing helper
│   └── redis.go                   # Redis client initialisation
├── internal/
│   ├── middleware/
│   │   └── rate_limiter.go        # IP-based rate limiting middleware
│   └── proxy/
│       ├── circuit_breaker.go     # Circuit breaker client wrapper
│       └── load_balancer.go       # Round-robin load balancer + reverse proxy
├── test/
│   ├── cb-test/
│   │   └── main.go                # Circuit breaker manual test
├── backend.go                     # Simulated backend service (run on multiple ports)
├── go.mod
└── go.sum
```

---

## 📦 Dependencies

| Library | Purpose |
|---|---|
| [`sony/gobreaker/v2`](https://github.com/sony/gobreaker) | Circuit breaker pattern implementation |
| [`redis/go-redis/v9`](https://github.com/redis/go-redis) | Redis client for rate limiting state |

---

## 🛠️ Prerequisites

- **Go** ≥ 1.21
- **Redis** running on `localhost:6379`

---

## ⚡ Quick Start

### 1. Clone the repo

```bash
git clone https://github.com/Mhabib34/go-api-gateway.git
cd go-api-gateway
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Start Redis

Make sure a Redis server is running locally:

```bash
redis-server
```

### 4. Start the backend services

Open two separate terminals and run a mock backend on each port:

```bash
# Terminal 1
go run backend.go -port 8081

# Terminal 2
go run backend.go -port 8082
```

### 5. Start the API Gateway

```bash
go run ./cmd/api
```

The gateway will start on port **8080**.

### 6. Test it

```bash
curl http://localhost:8080/
```

Responses will alternate between backend `8081` and `8082` thanks to round-robin load balancing.

---

## ⚙️ Configuration

All key settings are currently defined as constants or inline values. Here's a summary of what you can tune:

| Setting | Location | Default |
|---|---|---|
| Backend URLs | `main.go` slice | `localhost:8081`, `localhost:8082` | 
| Gateway port | `main.go` → `api` | `8080` |
| Rate limit — max requests | `rate_limiter.go` → `MaxRequest` | `100` |
| Rate limit — window | `rate_limiter.go` → `WindowTime` | `1 minute` |
| Circuit breaker — min requests to trip | `circuit_breaker.go` → `ReadyToTrip` | `5` |
| Circuit breaker — failure ratio to trip | `circuit_breaker.go` → `ReadyToTrip` | `60%` |
| Circuit breaker — half-open timeout | `circuit_breaker.go` → `Timeout` | `10 seconds` |
| Circuit breaker — max probe requests | `circuit_breaker.go` → `MaxRequests` | `3` |
| HTTP client timeout | `circuit_breaker.go` → `NewClient` | `3 seconds` |
| Redis address | `redis.go` | `localhost:6379` |

---

## 🔍 How It Works

### Request Flow

```
Client
  │
  ▼
Rate Limiter Middleware          ← checks IP-based quota via Redis
  │
  ▼
Load Balancer (Round-Robin)      ← picks the next alive backend
  │
  ▼
Reverse Proxy                    ← forwards request to chosen backend
  │
  ▼
Backend Service (8081 / 8082)
```

### Rate Limiter

Each incoming request is keyed by the client's IP address (extracted from `X-Forwarded-For` or `RemoteAddr`). A Redis `INCR` + `EXPIRE` pattern tracks how many requests have been made within the current one-minute window. Requests exceeding the limit receive a `429 Too Many Requests` response.

### Load Balancer

The load balancer maintains an ordered list of backends and cycles through them using a simple counter. Before forwarding, it checks the `Alive` flag on each backend. If a backend's reverse proxy returns an error, that backend is automatically marked as down and skipped on subsequent requests.

### Circuit Breaker

The circuit breaker (`sony/gobreaker`) wraps outbound HTTP calls made via the `Client` struct. It monitors success/failure ratios and, once the threshold is hit, opens the circuit — immediately failing requests without hitting the backend. After the timeout, it enters a half-open state and allows a limited number of probe requests to determine whether the backend has recovered.

---

## 🧪 Testing

### Circuit Breaker

There's a dedicated test script in `test/main.go` to manually verify the circuit breaker behaviour. It fires **10 consecutive requests** to an endpoint that always returns `500`, with a `500ms` delay between each — enough to clearly see the circuit trip in the logs.

Run it:

```bash
go run test/main.go
```

**What to expect in the logs:**

The first 5 requests will return a `server error: 500` because the circuit is still **closed** and counting failures. Once the failure ratio hits 60% on 5+ requests, the circuit **opens** — you'll see the state change log:

```
[CB] testing-service: closed → open
```

From that point on, remaining requests fail instantly with `gobreaker: open circuit` **without** actually hitting the endpoint. After 10 seconds of inactivity the circuit moves to **half-open** and allows up to 3 probe requests to test recovery.

---



To add more backends, simply update the `backends` slice in `main.go`:

```go
backends := []*proxy.Backend{
    {URL: config.MustParseURL("http://localhost:8081"), Alive: true},
    {URL: config.MustParseURL("http://localhost:8082"), Alive: true},
    {URL: config.MustParseURL("http://localhost:8083"), Alive: true}, // ← new backend
}
```

Then start a new backend process on the corresponding port:

```bash
go run cmd/backend/main.go -port 8083
```

---

## 📝 License

This project is open source. Feel free to use, modify, and contribute.
