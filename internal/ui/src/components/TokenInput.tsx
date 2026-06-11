import { useState, useEffect } from 'react'

interface TokenInputProps {
  onConnect: (token: string) => void
  defaultToken?: string
}

export default function TokenInput({ onConnect, defaultToken = '' }: TokenInputProps) {
  const [token, setToken] = useState(defaultToken)

  useEffect(() => {
    if (defaultToken && !token) {
      setToken(defaultToken)
    }
  }, [defaultToken])

  const handleConnect = () => {
    if (token.trim()) {
      onConnect(token)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleConnect()
    }
  }

  return (
    <div className="token-bar">
      <div className="token-inner">
        <label className="token-label" htmlFor="tok">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <rect x="3" y="11" width="18" height="11" rx="2" />
            <path d="M7 11V7a5 5 0 0110 0v4" />
          </svg>
          API Token
        </label>
        <div className="token-wrap">
          <input
            type="password"
            id="tok"
            placeholder="Enter your API token…"
            autoComplete="off"
            spellCheck="false"
            value={token}
            onChange={(e) => setToken(e.target.value)}
            onKeyPress={handleKeyPress}
          />
          <button className="btn btn-primary" onClick={handleConnect}>
            Connect
          </button>
        </div>
      </div>
    </div>
  )
}
