import { useCallback } from 'react'

interface FetchOptions extends RequestInit {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
}

export function useAPI(token?: string) {
  const fetch_api = useCallback(
    async <T>(url: string, options: FetchOptions = {}): Promise<T> => {
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...options.headers,
      }

      if (token) {
        headers['Authorization'] = token
      }

      const response = await fetch(url, {
        ...options,
        headers,
      })

      if (!response.ok) {
        const error = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(error.message || 'API request failed')
      }

      if (response.status === 204) {
        return undefined as T
      }

      return await response.json() as T
    },
    [token]
  )

  return { fetch_api }
}
