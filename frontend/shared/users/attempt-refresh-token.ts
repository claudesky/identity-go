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
