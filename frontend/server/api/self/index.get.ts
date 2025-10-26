import type { User } from '../../../app/composables/useAuth'
import type { DataResponse } from '../../../lib/response'

export default defineEventHandler(async (event) => {
  const accessToken = getCookie(event, 'identity_access_token')

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'No access token found',
    })
  }

  try {
    const response = await event.context.fetchWithAuth<DataResponse<User>>(
      '/self',
      { method: 'GET' }
    )

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Validation failed',
    })
  }
})
