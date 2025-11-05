import { H3Event } from 'h3'

export default defineEventHandler((event) => {
  // Skip interception for auth and non-api endpoints
  const path = event.path
  if (
    !path.startsWith('/api') ||
    path.startsWith('/api/auth/login') ||
    path.startsWith('/api/auth/register') ||
    path.startsWith('/api/auth/refresh')
  ) {
    return
  }

  // Add helper method to event context for retry logic
  event.context.fetchWithAuth = async <T>(
    url: string,
    options: any = {}
  ): Promise<T> => {
    try {
      return await apiFetch<T>(event, url, options)
    } catch (error: any) {
      // Retry on 401
      if (error.response?.status === 401) {
        const refreshed = await attemptTokenRefresh(event)

        if (refreshed) return await apiFetch<T>(event, url, options)
      }

      throw error
    }
  }
})

export function apiFetch<T>(event: H3Event, url: string, options: any) {
  const config = useRuntimeConfig()
  const accessToken = getCookie(event, 'identity_access_token')
  return $fetch<T>(`${config.public.apiBaseURL}${url}`, {
    ...options,
    headers: {
      ...options.headers,
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
    },
  })
}
