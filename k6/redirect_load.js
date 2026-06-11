import http from "k6/http"
import { sleep, check } from "k6"

const SLUG = "b5cfdab7"
const BASE_URL = "http://localhost:8080"

export const options = {
  stages: [
    { duration: "15s", target: 50 },   // ramp up to 50 users
    { duration: "30s", target: 200 },  // hold at 200 users
    { duration: "15s", target: 0 },    // ramp down
  ],
  thresholds: {
    http_req_duration: ["p(95)<200"],  // 95% of requests must finish under 200ms
    http_req_failed: ["rate<0.01"],    // less than 1% can fail
  },
}

export default function () {
  const res = http.get(`${BASE_URL}/${SLUG}`, {
    redirects: 0,
  })

  check(res, {
    "status is 302": (r) => r.status === 302,
    "has Location header": (r) => r.headers["Location"] !== undefined,
  })

  sleep(0.5)
}
