// @vitest-environment jsdom
import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, waitFor } from '@testing-library/react'
import { AuthProvider } from './AuthContext'
import { ResourceProvider } from './ResourceContext'

function Providers({ token }: { token: string }) {
  return (
    <AuthProvider token={token} onDisconnect={() => {}}>
      <ResourceProvider isConnected>{null}</ResourceProvider>
    </AuthProvider>
  )
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('ResourceProvider', () => {
  it('reloads resources when the connected token changes', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true, json: async () => ({}) })
    vi.stubGlobal('fetch', fetchMock)

    const { rerender } = render(<Providers token="first-token" />)
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(7))

    rerender(<Providers token="second-token" />)
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(14))

    expect(fetchMock.mock.calls.slice(7).every(([url, options]) =>
      (options as RequestInit).headers &&
      ((options as RequestInit).headers as Record<string, string>).Authorization === 'second-token',
    )).toBe(true)
  })
})
