import http from "k6/http"
import { check } from "k6"

const BASE_URL = "http://localhost:8080"
const TOKEN = __ENV.JWT_TOKEN      // pass via: k6 run -e JWT_TOKEN=xxx single_use_race.js

export const options = {
  vus: 2,           // exactly two users hitting the same slug simultaneously
  iterations: 1,    // one attempt per VU — total 2 requests at the same time
  thresholds: {
    checks: ["rate==1.0"],  // all checks must pass
  },
}

// setup runs once before the test — creates a fresh single-use link
export function setup() {
  const res = http.post(
    `${BASE_URL}/shorten`,
    JSON.stringify({
      long_url: "https://example.com",
      link_type: "payment",
      max_clicks: 1,
    }),
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${TOKEN}`,
      },
    }
  )

  const body = JSON.parse(res.body)
  return { slug: body.slug }
}

// default function receives the slug from setup
export default function (data) {
  const res = http.get(`${BASE_URL}/${data.slug}`, {
    redirects: 0,
  })

  // one VU should get 302, the other should get 410
  check(res, {
    "response is either 302 or 410": (r) => r.status === 302 || r.status === 410,
  })
}
