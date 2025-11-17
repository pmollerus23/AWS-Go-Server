import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'

interface HealthCheck {
  status: string
  timestamp?: string
}

function App() {
  const [count, setCount] = useState(0)
  const [health, setHealth] = useState<HealthCheck | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const checkHealth = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch('/healthz')
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      setHealth(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch health status')
      setHealth(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    checkHealth()
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
      <div className="max-w-4xl mx-auto px-8 py-16 text-center">
        {/* Logo Section */}
        <div className="flex justify-center items-center gap-8 mb-12">
          <a
            href="https://vite.dev"
            target="_blank"
            className="transition-transform hover:scale-110 hover:drop-shadow-[0_0_2em_rgba(100,108,255,0.6)]"
          >
            <img src={viteLogo} className="h-24 w-24" alt="Vite logo" />
          </a>
          <a
            href="https://react.dev"
            target="_blank"
            className="transition-transform hover:scale-110 hover:drop-shadow-[0_0_2em_rgba(97,218,251,0.6)] animate-[spin_20s_linear_infinite]"
          >
            <img src={reactLogo} className="h-24 w-24" alt="React logo" />
          </a>
        </div>

        {/* Title */}
        <h1 className="text-5xl font-bold mb-12 bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
          AWS Go Server + React
        </h1>

        {/* Health Check Card */}
        <div className="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 mb-8 border border-gray-700 shadow-xl">
          <h2 className="text-2xl font-semibold mb-6 text-gray-100">Backend Health Check</h2>
          <button
            onClick={checkHealth}
            disabled={loading}
            className="px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors duration-200 shadow-lg hover:shadow-blue-500/50"
          >
            {loading ? 'Checking...' : 'Check Backend Health'}
          </button>

          {health && (
            <div className="mt-6 p-4 bg-green-900/30 border border-green-700 rounded-lg">
              <p className="text-green-400 font-semibold">✓ Backend Status: {health.status}</p>
              {health.timestamp && (
                <p className="text-green-400/80 text-sm mt-1">Timestamp: {health.timestamp}</p>
              )}
            </div>
          )}

          {error && (
            <div className="mt-6 p-4 bg-red-900/30 border border-red-700 rounded-lg">
              <p className="text-red-400 font-semibold">✗ Error: {error}</p>
              <p className="text-red-400/80 text-sm mt-2">
                Make sure the Go server is running on port 8080
              </p>
            </div>
          )}
        </div>

        {/* Counter Card */}
        <div className="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 mb-8 border border-gray-700 shadow-xl">
          <button
            onClick={() => setCount((count) => count + 1)}
            className="px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg transition-colors duration-200 shadow-lg hover:shadow-purple-500/50"
          >
            count is {count}
          </button>
          <p className="mt-4 text-gray-400">
            Edit <code className="px-2 py-1 bg-gray-700 rounded text-sm text-blue-400">src/App.tsx</code> and save to test HMR
          </p>
        </div>

        {/* Footer Text */}
        <p className="text-gray-500 text-sm">
          This React app is integrated with your AWS Go Server backend
        </p>
      </div>
    </div>
  )
}

export default App
