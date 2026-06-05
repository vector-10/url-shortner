const BASE_URL = "http://localhost:8080"

function authHeaders(token: string) {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  }
}

export async function signup(email: string, password: string) {
  const res = await fetch(`${BASE_URL}/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function login(email: string, password: string) {
  const res = await fetch(`${BASE_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function shortenURL(longURL: string, token: string, linkType = "general", singleUse = false) {
  const res = await fetch(`${BASE_URL}/shorten`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({
      long_url: longURL,
      link_type: linkType,
      max_clicks: singleUse ? 1 : null,
    }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function listURLs(token: string) {
  const res = await fetch(`${BASE_URL}/urls`, {
    headers: authHeaders(token),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}
