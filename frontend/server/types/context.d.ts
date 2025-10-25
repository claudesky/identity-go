import type { FetchOptions } from 'ofetch'

declare module 'h3' {
  interface H3EventContext {
    fetchWithAuth: <T>(
      url: string,
      options?: FetchOptions
    ) => Promise<T>
  }
}

export {}
