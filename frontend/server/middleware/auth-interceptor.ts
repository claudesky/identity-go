/**
 * Server middleware to handle token refresh on 401 responses
 * Note: This middleware provides a helper but the actual retry logic
 * needs to be implemented per endpoint or using a fetch wrapper
 */
export default defineEventHandler((event) => {
  // Skip interception for auth endpoints
  const path = event.path
  if (
    path.includes('/api/auth/login') ||
    path.includes('/api/auth/register') ||
    path.includes('/api/auth/refresh')
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

function apiFetch<T>(event: any, url: string, options: any) {
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
