export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const body = await readBody(event)

  // Get the access token from cookie or Authorization header for authenticated requests
  let accessToken = getCookie(event, 'identity_access_token')
  if (!accessToken) {
    const authHeader = getHeader(event, 'authorization')
    if (authHeader && authHeader.startsWith('Bearer ')) {
      accessToken = authHeader.substring(7)
    }
  }

  // Forward all headers from the original request
  const headers = getHeaders(event)

  try {
    const response = await $fetch<any>(`${config.public.apiBaseURL}/auth/register`, {
      method: 'POST',
      body: body,
      headers: {
        ...headers,
        ...(accessToken ? { 'Authorization': `Bearer ${accessToken}` } : {})
      }
    })

    // Set cookies if tokens are present
    if (response) {
      if (response.access_token) {
        setCookie(event, 'identity_access_token', response.access_token, {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite: 'lax',
          maxAge: 60 * 15 // 15 minutes
        })
      }

      if (response.refresh_token) {
        setCookie(event, 'identity_refresh_token', response.refresh_token, {
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite: 'lax',
          maxAge: 60 * 60 * 24 * 7 // 7 days
        })
      }
    }

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Registration failed'
    })
  }
})
