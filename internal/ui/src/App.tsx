import { useState, useEffect } from 'react'
import { AuthProvider } from './context/AuthContext'
import { ResourceProvider } from './context/ResourceContext'
import Header from './components/Header'
import TokenInput from './components/TokenInput'
import MainView from './MainView'
import WelcomeScreen from './WelcomeScreen'

function App() {
  const [isConnected, setIsConnected] = useState(false)
  const [token, setToken] = useState<string>(() => {
    return localStorage.getItem('nullcloud_token') || ''
  })

  useEffect(() => {
    if (token) {
      setIsConnected(true)
    }
  }, [])

  const handleConnect = (authToken: string) => {
    setToken(authToken)
    localStorage.setItem('nullcloud_token', authToken)
    setIsConnected(true)
  }

  const handleDisconnect = () => {
    setToken('')
    localStorage.removeItem('nullcloud_token')
    setIsConnected(false)
  }

  return (
    <AuthProvider token={token} onDisconnect={handleDisconnect}>
      <ResourceProvider isConnected={isConnected}>
        <div className="app">
          <Header onDisconnect={handleDisconnect} />
          <TokenInput onConnect={handleConnect} defaultToken={token} />
          {!isConnected ? (
            <WelcomeScreen />
          ) : (
            <MainView />
          )}
        </div>
      </ResourceProvider>
    </AuthProvider>
  )
}

export default App
