import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"
import { listURLs } from "../api"
import { toast } from "sonner"

interface URLRecord {
  id: string
  slug: string
  long_url: string
  total_clicks: number
  expires_at: string | null
  created_at: string
  is_active: boolean
  max_clicks: number | null
  link_type: string
}

export default function Dashboard() {
  const [records, setRecords] = useState<URLRecord[]>([])
  const [loading, setLoading] = useState(true)

  const { token, logout } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) return
    listURLs(token)
      .then(setRecords)
      .catch(() => {
        toast.error("Could not load your links. Please try again.")
      })
      .finally(() => setLoading(false))
  }, [token])

  const handleLogout = () => {
    logout()
    navigate("/auth")
  }

  const handleQR = (slug: string) => {
    window.open(`http://localhost:8080/${slug}/qr`, "_blank")
  }

  const formatExpiry = (dateStr: string | null) => {
    if (!dateStr) return "Never"
    return new Date(dateStr).toLocaleString()
  }

  const getStatus = (record: URLRecord) => {
    if (!record.is_active) return { label: "Used", color: "text-zinc-500" }
    if (record.expires_at && new Date(record.expires_at) < new Date()) return { label: "Expired", color: "text-red-400" }
    return { label: "Active", color: "text-lime-400" }
  }

  return (
    <div className="min-h-screen bg-[#0a0a0a] flex flex-col">
      <header className="flex items-center justify-between px-8 py-5 border-b border-zinc-900">
        <span className="text-zinc-100 font-semibold tracking-tight">snip.</span>
        <div className="flex items-center gap-6">
          <button
            onClick={() => navigate("/")}
            className="text-xs text-zinc-500 hover:text-lime-400 transition-colors uppercase tracking-widest"
          >
            + New
          </button>
          <button
            onClick={handleLogout}
            className="text-xs text-zinc-500 hover:text-zinc-100 transition-colors uppercase tracking-widest"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="flex-1 px-8 py-8 max-w-6xl mx-auto w-full">
        <div className="flex items-end justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-zinc-100 tracking-tight">My Links</h1>
            <p className="text-zinc-600 text-sm mt-1">
              {records.length} link{records.length !== 1 ? "s" : ""} total
            </p>
          </div>
        </div>

        {loading && (
          <p className="text-zinc-600 text-sm">Loading your links...</p>
        )}

        {!loading && records.length === 0 && (
          <div className="border border-zinc-900 rounded-md p-16 text-center">
            <p className="text-zinc-600 text-sm mb-3">No links yet.</p>
            <button
              onClick={() => navigate("/")}
              className="text-lime-400 text-xs hover:underline uppercase tracking-widest"
            >
              Shorten your first URL
            </button>
          </div>
        )}

        {records.length > 0 && (
          <div className="border border-zinc-900 rounded-md overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-900">
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">Original</th>
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">Short link</th>
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">Type</th>
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">Clicks</th>
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">Status</th>
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">Expires</th>
                  <th className="px-4 py-3 text-left text-xs text-zinc-600 uppercase tracking-widest font-medium">QR</th>
                </tr>
              </thead>
              <tbody>
                {records.map((record, i) => {
                  const status = getStatus(record)
                  return (
                    <tr
                      key={record.id}
                      className={`border-b border-zinc-900 last:border-0 hover:bg-zinc-900 transition-colors ${i % 2 === 0 ? "bg-[#0a0a0a]" : "bg-[#0d0d0d]"}`}
                    >
                      <td className="px-4 py-3 max-w-xs">
                        <a
                          href={record.long_url}
                          target="_blank"
                          rel="noreferrer"
                          className="text-zinc-400 hover:text-zinc-100 transition-colors truncate block"
                        >
                          {record.long_url}
                        </a>
                      </td>
                      <td className="px-4 py-3">
                        <a
                          href={`http://localhost:8080/${record.slug}`}
                          target="_blank"
                          rel="noreferrer"
                          className="text-lime-400 font-mono text-xs hover:underline"
                        >
                          /{record.slug}
                        </a>
                      </td>
                      <td className="px-4 py-3">
                        <span className="text-zinc-500 text-xs uppercase tracking-widest">
                          {record.max_clicks === 1 ? "single-use" : record.link_type}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <span className="text-zinc-100 font-mono text-xs">{record.total_clicks}</span>
                      </td>
                      <td className="px-4 py-3">
                        <span className={`text-xs font-medium ${status.color}`}>{status.label}</span>
                      </td>
                      <td className="px-4 py-3 text-zinc-600 text-xs font-mono">
                        {formatExpiry(record.expires_at)}
                      </td>
                      <td className="px-4 py-3">
                        <button
                          onClick={() => handleQR(record.slug)}
                          className="text-zinc-600 hover:text-lime-400 transition-colors"
                          title="Open QR code"
                        >
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 9h6V3H3v6zm2-4h2v2H5V5zm8-2v6h6V3h-6zm4 4h-2V5h2v2zM3 21h6v-6H3v6zm2-4h2v2H5v-2zm13 4h2v-2h-2v2zm0-10h2v2h-2v-2zm-4 4h2v2h-2v-2zm2 4h-2v2h2v2h-2v-2h-2v2h-2v-4h4v-2h2v2zm-2-8h-2v2h2v-2z" />
                          </svg>
                        </button>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}
      </main>
    </div>
  )
}
