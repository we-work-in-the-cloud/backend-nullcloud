import React, { ReactNode } from 'react'

interface AuthContextType {
  token: string
  onDisconnect: () => void
}

export const AuthContext = React.createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ token, onDisconnect, children }: { token: string; onDisconnect: () => void; children: ReactNode }) {
  return (
    <AuthContext.Provider value={{ token, onDisconnect }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = React.useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
