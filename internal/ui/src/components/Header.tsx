import { useTheme } from '../hooks/useTheme'
import { useAuth } from '../context/AuthContext'

interface HeaderProps {
  onDisconnect: () => void
}

export default function Header({ onDisconnect }: HeaderProps) {
  const { theme, toggleTheme } = useTheme()
  const { token } = useAuth()

  return (
    <header>
      <a className="logo" href="/ui">
        <svg className="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
          <path d="M3 15a4 4 0 004 4h9a5 5 0 10-4.9-6H7a4 4 0 00-4 2z" />
        </svg>
        <span>NullCloud</span>
        <span className="logo-badge">Console</span>
      </a>
      <div className="header-right">
        {token && (
          <div className="conn-pill">
            <span className="conn-dot"></span>
            <span className="conn-token">{token.substring(0, 20)}...</span>
          </div>
        )}
        <button className="theme-toggle" onClick={toggleTheme} title="Toggle dark mode" aria-label="Toggle dark mode">
          <svg className="icon-moon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
          </svg>
          <svg className="icon-sun" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <circle cx="12" cy="12" r="5" />
            <line x1="12" y1="1" x2="12" y2="3" />
            <line x1="12" y1="21" x2="12" y2="23" />
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
            <line x1="1" y1="12" x2="3" y2="12" />
            <line x1="21" y1="12" x2="23" y2="12" />
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
          </svg>
        </button>
      </div>
    </header>
  )
}
