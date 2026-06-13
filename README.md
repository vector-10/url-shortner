# snip

Sending a link is simple. But in fintech, that link might be a payment request, an identity verification flow, or a new user's first step into your product. When it goes wrong — redeemed twice, expired too early, or never tracked — the consequence is not a broken page. It is a failed payment, a compliance gap, or a user who never completes onboarding.

snip is a link management service built around these use cases. Payment links expire in 48 hours and can only be redeemed once. KYC links give users 7 days. Onboarding links stay alive until acted on. Every redirect attempt is logged — because in regulated systems, "I think it worked" is not good enough.

The core engineering problem is atomicity. A single-use link that gets redeemed twice under concurrent load is not an edge case, it is a broken guarantee. snip solves this with a single atomic Postgres queries that makes the read and write indivisible.

---

## ERD

![snip ERD](snip-(url_shortener).png)

---

## What makes this different

**Single-use atomicity** — the naive implementation reads `is_active`, checks it, then sets it to false in a separate query. Two concurrent requests can both pass the check before either write lands. snip collapses the read and write into one atomic SQL statement. One request gets a 302, the other gets a 410. Postgres enforces this at the row level.

**Cache-aside with intentional invalidation** — Redis sits in front of Postgres, not as a primary store but as a read cache. On a single-use redemption, the cache entry is explicitly invalidated after the atomic write. A cache hit on a stale record cannot serve a dead link.

**Audit trail** — every redirect attempt is logged to `click_events` with its outcome and rejection reason. This is not optional in financial systems. If a payment link was used twice, you need to be able to prove when, from where, and why the second attempt was rejected or not.

**Per-type expiry** — link expiry is a business rule, not a technical default. Payment links expire in 48 hours. KYC links in 7 days. Onboarding links in 30 days. General links never expire unless they are single-use.

---

## Tech stack

| Layer | Choice | Why |
|-------|--------|-----|
| Backend | Go, `net/http` | Explicit, performant, no magic. Forces you to understand what the framework would otherwise hide. |
| Database | PostgreSQL | ACID guarantees. The atomic redemption pattern only works because Postgres provides row-level locking. |
| Cache | Redis | Sub-millisecond reads on the hot path. TTL-aware so cached entries never outlive their Postgres counterpart. |
| Auth | JWT + bcrypt + Google OAuth | Stateless tokens, no session store needed. OAuth handled via `golang.org/x/oauth2`. |
| Frontend | React, TypeScript, Tailwind CSS | Type-safe, dark theme, Sonner toasts. |
| Containers | Docker Compose | Backend, Postgres, and Redis wired together with healthchecks. |
| Load testing | k6 | Verified correctness and performance claims with real numbers. |

---

## Load test results

Two scenarios were tested against a single server instance with a 25-connection Postgres pool.

**Test 1 — Redirect throughput**

| Users | Duration | Requests | Throughput | Median | p95 | Error rate |
|-------|----------|----------|------------|--------|-----|------------|
| 200 | 60s | 8,617 | 143 req/s | 25ms | 693ms | 0% |
| 500 | 100s | 18,391 | 183 req/s | 525ms | 2.67s | 0% |

The gap between median and p95 at 500 users is Postgres connection pool queuing — the bottleneck at this scale is not the Go server or Redis, it is the 25-connection ceiling on the database. Zero errors across both runs.

**Test 2 — Single-use race condition**

Two virtual users hit the same single-use payment link simultaneously. One received a 302 redirect. The other received a 410 Gone. Checks passed 100%. The atomic write held under concurrent access.

---

## Running locally

**Prerequisites:** Docker, Docker Compose, Node.js (for the frontend)

```bash
# 1. Clone the repo
git clone https://github.com/vector-10/url-shortner.git
cd url-shortner

# 2. Set up environment variables
cp .env.example .env
# Fill in .env — see Environment Variables section below

# 3. Start backend, Postgres, and Redis
docker compose up --build -d

# 4. Start the frontend
cd client
npm install
npm run dev
```

Frontend runs at `http://localhost:5173`. Backend at `http://localhost:8080`.

---

## Environment variables

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `REDIS_ADDR` | Redis address — use `redis:6379` inside Docker |
| `POSTGRES_USER` | Postgres username |
| `POSTGRES_PASSWORD` | Postgres password |
| `POSTGRES_DB` | Postgres database name |
| `POSTGRES_URL` | Full Postgres connection string |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL — `http://localhost:8080/auth/google/callback` |
| `FRONTEND_URL` | Frontend origin for CORS and OAuth redirect — `http://localhost:5173` |
| `BASE_URL` | Backend base URL for QR code generation — `http://localhost:8080` |

Google OAuth requires a project in Google Cloud Console with the OAuth consent screen configured. Email/password auth works without it.

---

## What production would need

This is a single server demo. Getting to 10 million clicks per day honestly requires:

- **Horizontal scaling** — multiple Go instances behind a load balancer. The server is stateless so this is trivial.
- **PgBouncer** — connection pooling between the app and Postgres. At 10 instances × 25 connections you hit Postgres limits fast.
- **Redis Cluster** — the single Redis node is a single point of failure. A 3-primary cluster with replicas handles both availability and read distribution.
- **Async click logging** — writing to `click_events` on every redirect is 115 synchronous Postgres inserts per second at 10M clicks/day. A message queue draining into batch inserts removes this from the hot path.
- **CDN layer** — Cloudflare in front of the redirect endpoint caches responses at the edge. Popular links never reach the Go server.
