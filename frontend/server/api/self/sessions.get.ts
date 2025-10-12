import type { DataResponse } from '../../../lib/response'

export default defineEventHandler(async (event): Promise<any> => {
  const config = useRuntimeConfig()

  // Get the access token from cookie or Authorization header
  let accessToken = getCookie(event, 'identity_access_token')

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'Unauthorized',
    })
  }

  try {
    const response = await $fetch<DataResponse<any[]>>(
      `${config.public.apiBaseURL}/self/sessions`,
      {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    )

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Failed to fetch sessions',
    })
  }
})
