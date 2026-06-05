import { useState, type FormEvent } from "react"
import { toast } from "sonner"
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"
import { shortenURL } from "../api"

export default function Landing() {
  const [url, setUrl] = useState("")
  const [result, setResult] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)

  const { token } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!token) {
      navigate("/auth")
      return
    }
    setResult(null)
    setLoading(true)
    try {
      const data = await shortenURL(url, token)
      setResult(`http://localhost:8080/${data.slug}`)
      setUrl("")
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Something went wrong"

      if(message.includes("invalid") || message.includes("expired") || message.includes("token")) {
        toast.error("Session expired. Please sign in again.")
        setTimeout(() => navigate("/auth"), 1500)
      } else {
        toast.error(message)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCopy = () => {
    if (!result) return
    navigator.clipboard.writeText(result)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="min-h-screen bg-[#0a0a0a] flex flex-col">
      {/* Nav */}
      <nav className="flex items-center justify-between px-8 py-5 border-b border-zinc-900">
        <span className="text-zinc-100 font-semibold tracking-tight">snip.</span>
        <button
          onClick={() => navigate(token ? "/dashboard" : "/auth")}
          className="text-xs text-zinc-500 hover:text-zinc-100 transition-colors uppercase tracking-widest"
        >
          {token ? "Dashboard" : "Sign in"} →
        </button>
      </nav>

      {/* Main */}
      <main className="flex-1 flex flex-col items-center justify-center px-4">
        <div className="w-full max-w-2xl">
          <h1 className="text-6xl font-bold text-zinc-100 tracking-tighter mb-2 leading-tight">
            Long URLs,<br />
            <span className="text-lime-400">cut short.</span>
          </h1>
          <p className="text-zinc-500 text-base mb-10">
            Paste a URL. Get a clean link. No ads, no tracking.
          </p>

          <form onSubmit={handleSubmit} className="flex">
            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
              placeholder="https://example.com/your/very/long/url"
              className="flex-1 bg-zinc-900 text-zinc-100 text-sm px-4 py-3 border border-zinc-800 border-r-0 rounded-l-md outline-none placeholder:text-zinc-700 focus:border-zinc-600 transition-colors"
            />
            <button
              type="submit"
              disabled={loading}
              className="bg-lime-400 hover:bg-lime-300 disabled:opacity-40 text-zinc-900 font-semibold text-sm px-6 py-3 rounded-r-md transition-colors whitespace-nowrap"
            >
              {loading ? "..." : "Shorten"}
            </button>
          </form>

          {result && (
            <div className="mt-3 flex items-center justify-between border border-zinc-800 px-4 py-3 rounded-md bg-zinc-900">
              <a
                href={result}
                target="_blank"
                rel="noreferrer"
                className="text-lime-400 text-sm font-mono hover:underline truncate"
              >
                {result}
              </a>
              <button
                onClick={handleCopy}
                className="ml-4 text-xs text-zinc-500 hover:text-zinc-100 transition-colors whitespace-nowrap uppercase tracking-widest"
              >
                {copied ? "Copied" : "Copy"}
              </button>
            </div>
          )}
        </div>
      </main>

      {/* Footer */}
      <footer className="px-8 py-5 border-t border-zinc-900">
        <p className="text-zinc-700 text-xs">Links expire in 3 hours.</p>
      </footer>
    </div>
  )
}
