import { BrowserRouter, Routes, Route } from "react-router-dom"
import { AuthProvider } from "./context/AuthContext"
import ProtectedRoute from "./components/ProtectedRoute"
import Landing from "./pages/Landing"
import Auth from "./pages/Auth"
import { Toaster } from "sonner"
import Dashboard from "./pages/Dashboard"

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
      <Toaster richColors position="top-right" />      
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/auth" element={<Auth />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
        </Routes>        
      </BrowserRouter>
    </AuthProvider>
  )
}
