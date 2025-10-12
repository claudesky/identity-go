import type { User } from '../../../composables/useAuth'
import type { DataResponse } from '../../../lib/response'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()

  const accessToken = getCookie(event, 'identity_access_token')

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'No access token found',
    })
  }

  try {
    const response = await $fetch<DataResponse<User>>(
      `${config.public.apiBaseURL}/self`,
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
      message: error.data?.message || 'Validation failed',
    })
  }
})
