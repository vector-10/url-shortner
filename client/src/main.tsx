import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

const params = new URLSearchParams(window.location.search)
const token = params.get("token")
if (token) {
  localStorage.setItem("token", token)
  window.history.replaceState({}, "", "/dashboard")
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
