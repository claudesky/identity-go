export const useTokens = () => {
  const accessToken = useCookie('identity_access_token', {
    maxAge: 60 * 15, // 15 minutes
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    httpOnly: true
  })

  const refreshToken = useCookie('identity_refresh_token', {
    maxAge: 60 * 60 * 24 * 7, // 7 days
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    httpOnly: true
  })

  const clearTokens = () => {
    accessToken.value = null
    refreshToken.value = null
  }

  return {
    accessToken,
    refreshToken,
    clearTokens,
  }
}
