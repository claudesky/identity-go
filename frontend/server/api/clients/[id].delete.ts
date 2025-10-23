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

  const id = getRouterParam(event, 'id')

  if (!id) {
    throw createError({
      statusCode: 400,
      message: 'Client ID is required',
    })
  }

  try {
    const response = await $fetch<DataResponse<void>>(
      `${config.public.apiBaseURL}/clients/${id}`,
      {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    )

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Failed to delete client',
    })
  }
})
