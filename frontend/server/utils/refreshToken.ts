import type { H3Event } from 'h3'
import type { DataResponse, TokenData } from '../../lib/response'

/**
 * Attempt to refresh access token using refresh token from cookies
 * @param event H3Event to access cookies
 * @returns true if refresh succeeded, false otherwise
 */
export async function attemptTokenRefresh(event: H3Event): Promise<boolean> {
  const config = useRuntimeConfig()
  const refreshToken = getCookie(event, 'identity_refresh_token')

  if (!refreshToken) {
    return false
  }

  try {
    const response = await $fetch<DataResponse<TokenData>>(
      `${config.public.apiBaseURL}/auth/refresh`,
      {
        method: 'POST',
        body: {
          refresh_token: refreshToken,
        },
      }
    )

    setTokensFromData(event, response.data)

    return true
  } catch (error) {
    // Refresh failed, clear cookies
    deleteCookie(event, 'identity_access_token')
    deleteCookie(event, 'identity_refresh_token')
    return false
  }
}

export async function attemptRefreshToken(refreshToken: string) {
  const config = useRuntimeConfig()
  const response = await $fetch<DataResponse<TokenData>>(
    `${config.public.apiBaseURL}/auth/refresh`,
    {
      method: 'POST',
      body: { refresh_token: refreshToken },
    }
  )
  return response.data
}

export function setTokensFromData(event: H3Event, data: TokenData) {
  setCookie(event, 'identity_access_token', data.access_token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 15, // 15 minutes
  })

  setCookie(event, 'identity_refresh_token', data.refresh_token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 7, // 7 days
  })
}
