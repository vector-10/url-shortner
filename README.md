# url-shortner

A URL shortener REST API built with Go, using only the standard library and minimal dependencies. I built this as a learning project to get comfortable with Go's syntax, concurrency primitives, and interface-based design.
Syntax still feels weird but i'm getting there 😂

## Features

- Shorten any URL with an auto-generated or custom slug
- Redirect to original URL via slug
- Track click count per shortened URL
- View stats for any shortened URL
- Generate a QR code for any shortened URL
- TTL/expiry support with automatic background cleanup

## Tech

- Go 1.22+
- In-memory storage with `sync.RWMutex` for concurrent access
- Background goroutine for expired URL cleanup
- `net/http` ServeMux for routing (no external router)

## Project Structure

```
cmd/
  main.go              → entry point, wires everything together
internal/
  handler/
    handler.go         → HTTP handlers
  store/
    store.go           → Store interface
    memory.go          → in-memory implementation
  models/
    models.go          → URLRecord struct
```

## Running

```bash
go run cmd/main.go
```

Server starts on `http://localhost:8080`

## API

### Shorten a URL
```
POST /shorten
Content-Type: application/json

{ "long_url": "https://example.com" }
```

### Custom slug
```
POST /shorten
Content-Type: application/json

{ "long_url": "https://example.com", "slug": "myslug" }
```

### With expiry
```
POST /shorten
Content-Type: application/json

{ "long_url": "https://example.com", "expires_at": "2026-05-12T00:00:00Z" }
```

### Redirect
```
GET /:slug
```

### Stats
```
GET /:slug/stats
```

### QR Code
```
GET /:slug/qr
```
Returns a 256x256 PNG image of the QR code encoding the short URL.

## What's next

- Redis storage backend (drop-in via Store interface)
- Input validation
- Custom slug collision detection