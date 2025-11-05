import type { H3Event, EventHandlerRequest } from 'h3'

/**
 * Create authenticated fetch instance for making requests to the identity API
 * @param event H3Event for accessing cookies and context
 * @returns Fetch function with automatic auth token handling and refresh
 */
export function createIdentityApiFetch(event: H3Event<EventHandlerRequest>) {
  const config = useRuntimeConfig()

  return async <T = any>(url: string, options: any = {}): Promise<T> => {
    // NOTE: maybe use tokens from context in case another request authed?
    const accessToken = getCookie(event, 'identity_access_token')
    const refreshToken = getCookie(event, 'identity_refresh_token')

    // Preemptively refresh if no access token but have refresh token
    if (!accessToken && refreshToken) {
      await attemptTokenRefresh(event)
    }

    // Get token (from context if just refreshed, or from cookie)
    const token = event.context.accessToken ?? getCookie(event, 'identity_access_token')

    // Make initial request
    try {
      return await $fetch<T>(`${config.public.apiBaseURL}${url}`, {
        ...options,
        headers: {
          ...options.headers,
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      })
    } catch (error: any) {
      // Retry on 401 after attempting token refresh
      if (error.response?.status === 401) {
        const refreshed = await attemptTokenRefresh(event)
        if (refreshed) {
          const newToken = event.context.accessToken ?? getCookie(event, 'identity_access_token')
          return $fetch<T>(`${config.public.apiBaseURL}${url}`, {
            ...options,
            headers: {
              ...options.headers,
              ...(newToken ? { Authorization: `Bearer ${newToken}` } : {}),
            },
          })
        }
      }

      throw error
    }
  }
}
