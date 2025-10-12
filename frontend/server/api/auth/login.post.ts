import type { DataResponse, TokenData } from "../../../lib/response"

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const body = await readBody(event)

  try {
    const response = await $fetch<DataResponse<TokenData>>(
      `${config.public.apiBaseURL}/auth/login`,
      {
        method: 'POST',
        body: body,
      }
    )

    setCookie(event, 'identity_access_token', response.data.access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 60 * 15, // 15 minutes
    })

    setCookie(event, 'identity_refresh_token', response.data.refresh_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 60 * 60 * 24 * 7, // 7 days
    })

    // Return success without exposing tokens
    return
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Login failed',
    })
  }
})
