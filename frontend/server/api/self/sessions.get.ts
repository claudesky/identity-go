import type { DataResponse } from '../../../lib/response'

export default defineEventHandler(async (event): Promise<any> => {
  const accessToken = getCookie(event, 'identity_access_token')

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'Unauthorized',
    })
  }

  try {
    const response = await event.context.fetchWithAuth<DataResponse<any[]>>(
      '/self/sessions',
      { method: 'GET' }
    )

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Failed to fetch sessions',
    })
  }
})
